package server

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/TwoWaySix/ytdl-gui/internal"
)

//go:embed template/index.html
var indexTemplate string

//go:embed static/script.js
var javascript string

//go:embed static/style.css
var css string

type TemplateData struct {
	CSS    template.CSS
	JS     template.JS
	Status string
}

type YouTubeDownloadServer struct {
	address       string
	mux           http.ServeMux
	status        string
	downloadItems map[string]DownloadItem
}

type DownloadItem struct {
	Url    string
	Title  string
	Status string
}

func NewYouTubeDownloadServer(address string) YouTubeDownloadServer {
	backend := YouTubeDownloadServer{
		address:       address,
		mux:           *http.NewServeMux(),
		status:        "Everything normal",
		downloadItems: make(map[string]DownloadItem, 0),
	}
	backend.mux.HandleFunc("/", backend.ServeIndex)
	backend.mux.HandleFunc("/download", backend.HandleDownloadRequest)
	backend.mux.HandleFunc("/status", backend.GetServerStatus)
	backend.mux.HandleFunc("/items", backend.GetDownloadItems)
	return backend
}

func (ys *YouTubeDownloadServer) Start() error {
	log.Printf("INFO: Visit yt-dlp GUI in your web browser via http://%s", ys.address)
	err := http.ListenAndServe(ys.address, &ys.mux)
	if err != nil {
		return fmt.Errorf("starting server: %w", err)
	}
	return nil
}

func (ys *YouTubeDownloadServer) ServeIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	data := TemplateData{
		CSS:    template.CSS(css),
		JS:     template.JS(javascript),
		Status: ys.status,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (ys *YouTubeDownloadServer) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(map[string]interface{}{"message": ys.status})
	if err != nil {
		ys.logAndSetError(fmt.Errorf("getting message: %w", err))
		return
	}
	w.Write(j)
}

// TODO: Display download items in frontend
func (ys *YouTubeDownloadServer) GetDownloadItems(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(map[string]interface{}{"items": ys.downloadItems})
	if err != nil {
		ys.logAndSetError(fmt.Errorf("getting download items: %w", err))
		return
	}
	w.Write(j)
}

func (ys *YouTubeDownloadServer) HandleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	videoUrl := r.URL.Query().Get("url")
	if videoUrl == "" {
		ys.logAndSetError(errors.New("empty video url provided"))
		return
	}

	videoTitle, err := internal.ExtractWebsiteTitle(videoUrl)
	if err != nil {
		ys.logAndSetError(fmt.Errorf("handling download: %w", err))
		return
	}

	item := DownloadItem{
		Url:    videoUrl,
		Title:  videoTitle,
		Status: "Downloading",
	}
	ys.downloadItems[videoUrl] = item

	ys.status = fmt.Sprintf("Started downloading '%s' from '%s'", videoTitle, videoUrl)
	err = internal.DownloadVideo(videoUrl)
	if err != nil {
		ys.logAndSetError(fmt.Errorf("handling download: %w", err))
		return
	}
	item.Status = "Finished"
	ys.downloadItems[videoUrl] = item

	ys.status = fmt.Sprintf("Finished downloading '%s' from '%s'", videoTitle, videoUrl)
}

func (ys *YouTubeDownloadServer) logAndSetError(err error) {
	log.Printf("ERROR: %v", err)
	ys.status = fmt.Sprintf("ERROR: %v", err)
}
