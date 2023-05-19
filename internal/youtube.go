package internal

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

var YouTubeDownloaderExeURL = "https://yt-dl.org/latest/youtube-dl.exe"

func DownloadVideo(url string) error {
	log.Printf("Trying to extract audio from '%s'", url)
	var err error
	if runtime.GOOS == "windows" {
		cmd := exec.Command("./youtube-dl", "-x", "--audio-format", "mp3", url)
		err = cmd.Run()
	} else {
		cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", url)
		err = cmd.Run()
	}

	return err
}

func AssureYouTubeDownloader() error {
	var err error
	if runtime.GOOS == "windows" {
		err = DownloadYtdlExe()
	} else {
		_, err = exec.LookPath("youtube-dl")
	}
	return err
}

func DownloadYtdlExe() error {
	log.Printf("INFO: Downloading '%s'", YouTubeDownloaderExeURL)
	response, err := http.Get(YouTubeDownloaderExeURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	f, err := os.Create("youtube-dl.exe")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return err
	}
	return nil
}
