package main

import (
	"github.com/usmon1983/http/pkg/server"
	"log"
	"net"
	"os"
)

func main() {
	host := "0.0.0.0"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}

func execute(host string, port string) (err error)  {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", func (req *server.Request)  {
		id := req.QueryParams["id"]
		log.Println(id)
	})
	return srv.Start()
}