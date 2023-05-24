package main

import (
	"fmt"
	"os"

	"github.com/TwoWaySix/ytdl-gui/internal"
	"github.com/TwoWaySix/ytdl-gui/internal/server"
)

var address = "127.0.0.1:2345"
var downloadDirectory = "downloads"

func main() {
	err := internal.AssurePrerequisites()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(-1)
	}

	backend := server.NewYouTubeDownloadServer(address, downloadDirectory)
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
