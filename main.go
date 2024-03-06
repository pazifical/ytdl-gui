package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/pazifical/ytdl-gui/internal"
	"github.com/pazifical/ytdl-gui/internal/server"
)

//go:embed templates
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

var address = "127.0.0.1:2345"
var downloadDirectory = "downloads"

func main() {
	err := internal.AssurePrerequisites(downloadDirectory)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(-1)
	}

	backend, err := server.NewYouTubeDownloadServer(address, downloadDirectory, templateFS, staticFS)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(-1)
	}

	err = backend.Start()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(-1)
	}
}
