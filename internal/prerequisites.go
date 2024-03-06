package internal

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

func AssurePrerequisites(downloadDirectory string) error {
	err := assureYouTubeDownloader()
	if err != nil {
		return fmt.Errorf("assuring prerequisites: %w", err)
	}
	err = assureFfmpeg()
	if err != nil {
		return fmt.Errorf("assuring prerequisites: %w", err)
	}

	err = assureDownloadDirectory(downloadDirectory)
	if err != nil {
		return fmt.Errorf("assuring prerequisites: %w", err)
	}
	return nil
}

func assureDownloadDirectory(directory string) error {
	_, err := os.Stat(directory)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(directory, 0750)
		if err != nil {
			return err
		}
	} else if err != nil {
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
	if err != nil {
		return fmt.Errorf("assuring yt-dlp: %w", err)
	}
	return err
}

func assureFfmpeg() error {
	ffmpegZipPath := "ffmpeg.zip"

	var err error
	if runtime.GOOS == "windows" {
		// TODO: Properly implement
		return nil
		err = downloadAndExtractFfmpeg(ffmpegZipPath)
	} else {
		_, err = exec.LookPath("ffmpeg")
	}
	if err != nil {
		return fmt.Errorf("assuring ffmpeg: %w", err)
	}
	return err
}

func downloadYtdlExe() error {
	errorTemplate := "downloading yt-dlp: %w"
	log.Printf("INFO: Downloading '%s'", youTubeDownloaderExeURL)
	response, err := http.Get(youTubeDownloaderExeURL)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	f, err := os.Create("yt-dlp.exe")
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}

func downloadAndExtractFfmpeg(ffmpegZipPath string) error {
	errorTemplate := "downloading and extracting ffmpeg: %w"

	err := downloadFfmpeg(ffmpegZipPath)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	err = extractFfmpeg(ffmpegZipPath)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}

	return nil
}

func downloadFfmpeg(ffmpegZipPath string) error {
	errorTemplate := "downloading ffmpeg: %w"

	log.Printf("INFO: Downloading '%s'", ffmpegDownloaderExeURL)
	response, err := http.Get(youTubeDownloaderExeURL)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	defer response.Body.Close()

	f, err := os.Create(ffmpegZipPath)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	defer f.Close()

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}

func extractFfmpeg(ffmpegZipPath string) error {
	archive, err := zip.OpenReader(ffmpegZipPath)
	if err != nil {
		return fmt.Errorf("downloading and extracting ffmpeg: %w", err)
	}
	defer archive.Close()
	for _, f := range archive.File {
		fmt.Println(f.Name)
		if f.FileInfo().IsDir() {
			continue
		}
		// TODO: Implement
	}
	return nil
}
