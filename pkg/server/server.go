package server

import (
	"net/url"
	"strings"
	"bytes"
	"io"
	"log"
	"sync"
	"net"
)

type HandlerFunc func (req *Request)

type Server struct {
	addr string
	mu sync.RWMutex
	handlers map[string]HandlerFunc
}

type Request struct {
	Conn net.Conn
	QueryParams url.Values
}

func NewServer(addr string) *Server  {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc)  {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error  {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func ()  {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()		

	for {
	conn, err := listener.Accept()
		
		if err != nil {
			log.Print(err)
			//идем обслуживать следующего
			continue
		}
		go s.handle(conn)
				
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, (1024 * 8))
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Printf("%s", buf[:n])
		}
		if err != nil {
			log.Println(err)
			return
		}

		var req Request
		data := buf[:n]
		requestLineDelim := []byte{'\r', '\n'}
		requestLineEnd := bytes.Index(data, requestLineDelim)
		if requestLineEnd == -1 {
			log.Printf("Bad Request")
			return
		}

		requestLine := string(data[:requestLineEnd])
		parts := strings.Split(requestLine, " ")

		if len(parts) != 3 {
			return
		}
		path, version := parts[1], parts[2]
		if version != "HTTP/1.1" {
			return
		}

		decode, err := url.PathUnescape(path)
		if err != nil {
			log.Println(err)
			return
		}

		uri, err := url.ParseRequestURI(decode)
		if err != nil {
			log.Println(err)
			return
		}
		
		req.Conn = conn
		req.QueryParams = uri.Query()
		
		var handler = func(req *Request) { conn.Close() }

		s.mu.RLock()
		for i := 0; i < len(s.handlers); i++ {
			if hr, found := s.handlers[uri.Path]; found {
				handler = hr
				break
			}
		}
		s.mu.RUnlock()

		handler(&req)
	}
}