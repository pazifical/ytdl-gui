package main

import (
	"github.com/TwoWaySix/ytdl-gui/internal"
	"github.com/TwoWaySix/ytdl-gui/internal/server"
)

var address = "127.0.0.1:3000"

func main() {
	err := internal.AssurePrerequisites()
	if err != nil {
		panic(err)
	}

	backend := server.NewYtdlpServer(address)
	err = backend.Start()
	if err != nil {
		panic(err)
	}

	// http.HandleFunc("/", server.ServeIndex)
	// http.HandleFunc("/download", server.HandleDownload)

	// log.Printf("INFO: Visit yt-dlp GUI in your web browser via http://%s", address)
	// err = http.ListenAndServe(address, nil)
	// if err != nil {
	// 	panic(err)
	// }
}
