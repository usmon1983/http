package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"github.com/usmon1983/http/pkg/banners"
	"github.com/usmon1983/http/cmd/app"
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

func execute(host string, port string) (err error) {
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	server := app.NewServer(mux, bannersSvc)
	
	server.Init()
	srv := &http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: http.HandlerFunc(func (writer http.ResponseWriter, request *http.Request)  {
			log.Print(request.RequestURI) //полный URI
			log.Print(request.Method) //метод
			//log.Print(request.Header) //все заголовки
			//log.Print(request.Header.Get("Content-Type")) //конкретный заголовок

			//log.Print(request.FormValue("tags")) //только первое значение Query + POST
			//log.Print(request.PostFormValue("tags")) //только первое значение POST

			body, err := ioutil.ReadAll(request.Body) //тело запроса
			if err != nil {
				log.Print(err)
			}
			log.Print("%s", body)

			err = request.ParseMultipartForm(10 * 1024 * 1024) //10Mb
			if err != nil {
				log.Print(err)
			}
			
			//доступно только после ParseForm (либо FormValue, PostFormValue)
			//log.Print(request.Form)
			//log.Print(request.PostForm)
			//доступно только после ParseMultipart (FormValue, PostFormValue автоматически вызывают ParseMultipartForm)
			//log.Print(request.FormFile("image"))
			//request.MultipartForm.Value - только "обычные поля"
			//request.MultipartForm.File - только файлы
			
		}),
	}
	return srv.ListenAndServe()
}