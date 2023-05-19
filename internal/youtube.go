package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		cmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", url)
		err = cmd.Run()
	}

	log.Printf("Finished extracting audio from '%s'", url)
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

func AssureFfmpeg() error {
	var err error
	if runtime.GOOS == "windows" {
		err = DownloadAndExtractFfmpeg()
	} else {
		_, err = exec.LookPath("ffmpeg")
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

	f, err := os.Create("yt-dlp.exe")
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

func DownloadAndExtractFfmpeg() error {
	log.Printf("INFO: Downloading '%s'", FfmpegDownloaderExeURL)
	response, err := http.Get(YouTubeDownloaderExeURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	f, err := os.Create("ffmpeg.zip")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return err
	}

	// TODO: Extract from zip
	archive, err := zip.OpenReader("ffmpeg.zip")
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, f := range archive.File {
		fmt.Println(f.Name)
		if f.FileInfo().IsDir() {
			continue
		}
	}

	return nil
}
