package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pazifical/ytdl-gui/internal"
)

func (ys *YouTubeDownloadServer) ServeIndex(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Status: ys.status,
	}

	err := ys.templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (ys *YouTubeDownloadServer) GetServerStatus(w http.ResponseWriter, r *http.Request) {
	err := ys.templates.ExecuteTemplate(w, "status.html", ys.status)
	if err != nil {
		ys.logAndSetError(fmt.Errorf("getting download items: %w", err))
	}
}

// TODO: Display download items in frontend
func (ys *YouTubeDownloadServer) GetDownloadItems(w http.ResponseWriter, r *http.Request) {
	dirEntries, err := os.ReadDir(ys.downloadDirectory)
	if err != nil {
		ys.logAndSetError(fmt.Errorf("getting download items: %w", err))
		return
	}

	downloads := make([]DownloadItem, 0)
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}
		name := dirEntry.Name()
		downloads = append(downloads, DownloadItem{
			Title:  name,
			Status: "Finished",
		})
	}

	err = ys.templates.ExecuteTemplate(w, "downloads.html", downloads)
	if err != nil {
		ys.logAndSetError(fmt.Errorf("getting download items: %w", err))
		return
	}
}

func (ys *YouTubeDownloadServer) HandleDownloadRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	videoUrl := r.Form.Get("url")
	if videoUrl == "" {
		log.Printf("ERROR: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	go internal.DownloadVideo(videoUrl, ys.downloadDirectory)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (ys *YouTubeDownloadServer) logAndSetError(err error) {
	log.Printf("ERROR: %v", err)
	ys.status = fmt.Sprintf("ERROR: %v", err)
}
