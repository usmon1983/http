package app

import (
	"strings"
	"encoding/json"
	"log"
	"strconv"
	"github.com/usmon1983/http/pkg/banners"
	"net/http"
)

//Server представляет собой логический сервер нашего приложения
type Server struct {
	mux *http.ServeMux
	bannersSvc *banners.Service
}

//NewServer - функция конструктор для создания сервера
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writter http.ResponseWriter, request *http.Request)  {
	s.mux.ServeHTTP(writter, request)
}

//Init инициализирует сервер (регистрирует все Handler'ы)
func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request)  {
	banners, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request)  {
	idParam := request.URL.Query().Get("id")
	
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}


func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request)  {
	idParam := request.URL.Query().Get("id")
	title := request.URL.Query().Get("title")
	content := request.URL.Query().Get("content")
	button := request.URL.Query().Get("button")
	link := request.URL.Query().Get("link")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if title == "" && content == "" && button == "" && link == "" {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	item := &banners.Banner{
		ID:      id,
		Title:   title,
		Content: content,
		Button:  button,
		Link:    link,
	}

	file, header, err := request.FormFile("image")
	if err == nil { 
		var name = strings.Split(header.Filename, ".")
		item.Image = name[len(name)-1]
	}

	banner, err := s.bannersSvc.Save(request.Context(), item, file)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}	
}

func (s *Server) handleRemoveById(writer http.ResponseWriter, request *http.Request)  {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	banner, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}