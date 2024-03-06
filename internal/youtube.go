package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var youTubeDownloaderExeURL = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"

var ffmpegDownloaderExeURL = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-git-full.7z"

// TODO: Download into subdirectory
func DownloadVideo(url string, directory string) error {
	log.Printf("Trying to extract audio from '%s'", url)

	errorTemplate := "downloading video: %w"

	var tool string
	if runtime.GOOS == "windows" {
		tool = "./yt-dlp"
	} else {
		tool = "yt-dlp"
	}

	var builder strings.Builder
	builder.WriteString(directory)
	builder.WriteString("/")
	builder.WriteString("%(title)s.%(ext)s")

	cmd := exec.Command(tool, "-o", builder.String(), "-x", "--audio-format", "mp3", url)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	errorLogs, err := io.ReadAll(stderr)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	if string(errorLogs) != "" {
		return fmt.Errorf("downloading video: %s", errorLogs)
	}

	log.Printf("Finished extracting audio from '%s'", url)
	return err
}

func ExtractWebsiteTitle(url string) (string, error) {
	errorTemplate := "extracting title: %w"
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf(errorTemplate, err)
	}
	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf(errorTemplate, err)
	}

	content := string(bytes)
	content = strings.ReplaceAll(content, "\n", "")

	r, err := regexp.Compile("<title>.*</title>")
	if err != nil {
		return "", fmt.Errorf(errorTemplate, err)
	}

	match := r.FindString(content)
	match = strings.ReplaceAll(match, "<title>", "")
	match = strings.ReplaceAll(match, "</title>", "")

	return match, nil
}
