package server

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	Status string
}

type YouTubeDownloadServer struct {
	address           string
	mux               http.ServeMux
	status            string
	downloadDirectory string
	templates         *template.Template
}

type DownloadItem struct {
	Url    string
	Title  string
	Status string
}

func NewYouTubeDownloadServer(address, downloadDirectory string, frontendFS, staticFS embed.FS) (YouTubeDownloadServer, error) {
	tmpls, err := template.ParseFS(frontendFS, "templates/*.html")
	if err != nil {
		return YouTubeDownloadServer{}, err
	}

	backend := YouTubeDownloadServer{
		address:           address,
		mux:               *http.NewServeMux(),
		status:            "Everything normal",
		downloadDirectory: downloadDirectory,
		templates:         tmpls,
	}

	backend.mux.HandleFunc("/", backend.ServeIndex)
	backend.mux.HandleFunc("POST /api/download", backend.HandleDownloadRequest)
	backend.mux.HandleFunc("GET /api/html/status", backend.GetServerStatus)
	backend.mux.HandleFunc("GET /api/html/downloads", backend.GetDownloadItems)
	backend.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	return backend, nil
}

func (ys *YouTubeDownloadServer) Start() error {
	errorTemplate := "starting server: %w"

	log.Printf("INFO: Visit yt-dlp GUI in your web browser via http://%s", ys.address)
	err := http.ListenAndServe(ys.address, &ys.mux)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}
