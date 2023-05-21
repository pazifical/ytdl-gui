package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var YouTubeDownloaderExeURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"

var FfmpegDownloaderExeURL = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-git-full.7z"

func DownloadVideo(url string) error {
	log.Printf("Trying to extract audio from '%s'", url)
	var tool string
	if runtime.GOOS == "windows" {
		tool = "./yt-dlp"
	} else {
		tool = "yt-dlp"
	}

	cmd := exec.Command(tool, "-x", "--audio-format", "mp3", url)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	slurp, err := io.ReadAll(stderr)
	if err != nil {
		return err
	}

	if string(slurp) != "" {
		return errors.New(string(slurp))
	}

	log.Printf("Finished extracting audio from '%s'", url)
	return err
}

func ExtractTitle(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	content := string(bytes)
	content = strings.ReplaceAll(content, "\n", "")

	r, err := regexp.Compile("<title>.*</title>")
	if err != nil {
		return "", err
	}

	match := r.FindString(content)
	match = strings.ReplaceAll(match, "<title>", "")
	match = strings.ReplaceAll(match, "</title>", "")
	return match, nil
}
