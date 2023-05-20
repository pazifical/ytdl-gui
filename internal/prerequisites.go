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

func AssurePrerequisites() error {
	err := assureYouTubeDownloader()
	if err != nil {
		return err
	}
	err = assureFfmpeg()
	if err != nil {
		return err
	}
	return nil
}

func assureYouTubeDownloader() error {
	var err error
	if runtime.GOOS == "windows" {
		err = downloadYtdlExe()
	} else {
		_, err = exec.LookPath("yt-dlp")
	}
	return err
}

func assureFfmpeg() error {
	var err error
	if runtime.GOOS == "windows" {
		err = downloadAndExtractFfmpeg()
	} else {
		_, err = exec.LookPath("ffmpeg")
	}
	return err
}

func downloadYtdlExe() error {
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

func downloadAndExtractFfmpeg() error {
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
