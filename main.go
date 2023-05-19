package main

import (
	_ "embed"
	"log"
	"net/http"
	"text/template"

	"github.com/TwoWaySix/ytdl-gui/internal"
)

//go:embed template/index.html
var indexTemplate string

//go:embed static/script.js
var javascript string

//go:embed static/style.css
var css string

func main() {
	err := internal.AssureYouTubeDownloader()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/download", handleDownload)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = t.Execute(w, struct{ Css, JS string }{css, javascript})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	videoUrl := r.URL.Query().Get("url")
	if videoUrl == "" {
		panic("Empty video url provided")
	}
	err := internal.DownloadVideo(videoUrl)
	if err != nil {
		panic(err)
	}
}
