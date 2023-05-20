package internal

import (
	"log"
	"os/exec"
	"runtime"
)

var YouTubeDownloaderExeURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"

var FfmpegDownloaderExeURL = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-git-full.7z"

func DownloadVideo(url string) error {
	log.Printf("Trying to extract audio from '%s'", url)
	var err error
	if runtime.GOOS == "windows" {
		cmd := exec.Command("./yt-dlp", "-x", "--audio-format", "mp3", url)
		err = cmd.Run()
	} else {
		cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", url)
		err = cmd.Run()
	}

	log.Printf("Finished extracting audio from '%s'", url)
	return err
}
