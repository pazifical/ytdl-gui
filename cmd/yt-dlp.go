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
}
