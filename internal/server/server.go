package server

import (
	_ "embed"
	"encoding/json"
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

type Data struct {
	CSS     template.CSS
	JS      template.JS
	Message string
}

type YtdlpServer struct {
	address string
	mux     http.ServeMux
	Message string
}

func NewYtdlpServer(address string) YtdlpServer {
	backend := YtdlpServer{address: address, mux: *http.NewServeMux(), Message: "Everything normal"}
	backend.mux.HandleFunc("/", backend.ServeIndex)
	backend.mux.HandleFunc("/download", backend.HandleDownload)
	backend.mux.HandleFunc("/status", backend.GetMessage)
	return backend
}

func (ys *YtdlpServer) Start() error {
	log.Printf("INFO: Visit yt-dlp GUI in your web browser via http://%s", ys.address)
	err := http.ListenAndServe(ys.address, &ys.mux)
	if err != nil {
		return err
	}
	return nil
}

func (ys *YtdlpServer) ServeIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	data := Data{
		CSS:     template.CSS(css),
		JS:      template.JS(javascript),
		Message: ys.Message,
	}

	err = t.Execute(w, data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (ys *YtdlpServer) GetMessage(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(map[string]interface{}{"message": ys.Message})
	w.Write(j)
}

func (ys *YtdlpServer) HandleDownload(w http.ResponseWriter, r *http.Request) {
	videoUrl := r.URL.Query().Get("url")
	if videoUrl == "" {
		log.Println("Empty video url provided")
		ys.Message = "Please enter an URL"
		return
	}

	videoTitle, err := internal.ExtractTitle(videoUrl)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ys.Message = fmt.Sprintf("ERROR: %v", err)
		return
	}

	ys.Message = fmt.Sprintf("Started downloading '%s' from '%s'", videoTitle, videoUrl)
	err = internal.DownloadVideo(videoUrl)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ys.Message = fmt.Sprintf("ERROR: %v", err)
		return
	}
	ys.Message = fmt.Sprintf("Finished downloading '%s' from '%s'", videoTitle, videoUrl)
}
