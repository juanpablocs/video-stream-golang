package usecases

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

type Resolution struct {
	Width  int
	Height int
	Label  string // Como "720p", "480p", etc.
}

func TranscodeVideo(originalInputPath, outputPath string, totalDuration time.Duration) error {
	currentInputPath := originalInputPath // Empieza con el path original

	resolutions := []Resolution{
		{Width: 1280, Height: 720, Label: "720p"},
		{Width: 854, Height: 480, Label: "480p"},
		{Width: 640, Height: 360, Label: "360p"},
		{Width: 426, Height: 240, Label: "240p"},
	}

	for i, res := range resolutions {
		fmt.Printf("\nIniciando transcodificación a %s...\n", res.Label)

		// Definir los nombres de los archivos de salida basados en la resolución actual
		basePath := path.Dir(outputPath)

		parts := strings.Split(outputPath, "/")
		name := parts[len(parts)-1]
		resOutputMP4 := fmt.Sprintf("%s/%s_%s.mp4", basePath, name, res.Label)
		resOutputHLS := fmt.Sprintf("%s/%s_%s.m3u8", basePath, name, res.Label)

		// Crear las carpetas necesarias
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("error al crear las carpetas necesarias: %w", err)
		}

		// Ejecutar la transcodificación para la resolución actual
		if err := executeAndMonitor(currentInputPath, resOutputMP4, resOutputHLS, res); err != nil {
			return err
		}

		// Para la siguiente iteración, usar el MP4 que acabamos de crear como input
		currentInputPath = resOutputMP4

		// Para evitar re-codificar el original en la primera iteración, restaurar el inputPath original
		if i == 0 {
			currentInputPath = originalInputPath
		}
	}

	fmt.Println("\nProcesos completados con éxito.")
	return nil
}

func executeAndMonitor(inputPath, outputPathMP4, outputPathHLS string, res Resolution) error {
	var wg sync.WaitGroup
	wg.Add(1)

	// Transcodificar a MP4
	fmt.Printf("Transcodificando a MP4: %s\n", res.Label)
	cmdMP4 := exec.Command("ffmpeg", "-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:%d", res.Width, res.Height), "-c:v", "libx264",
		"-preset", "slow", "-crf", "23",
		"-c:a", "aac", "-b:a", "128k",
		outputPathMP4)

	if err := cmdMP4.Start(); err != nil {
		return fmt.Errorf("error al iniciar MP4 Transcoding %s: %w", res.Label, err)
	}

	go func() {
		defer wg.Done()
	}()

	wg.Wait()

	if err := cmdMP4.Wait(); err != nil {
		return fmt.Errorf("MP4 Transcoding %s terminó con error: %w", res.Label, err)
	}

	// Generar HLS
	fmt.Printf("Generando HLS: %s\n", res.Label)
	cmdHLS := exec.Command("ffmpeg", "-i", outputPathMP4,
		"-hls_time", "10", "-hls_playlist_type", "vod",
		"-c:v", "copy", "-c:a", "copy",
		outputPathHLS)

	if err := cmdHLS.Start(); err != nil {
		return fmt.Errorf("error al iniciar HLS Generation %s: %w", res.Label, err)
	}

	if err := cmdHLS.Wait(); err != nil {
		return fmt.Errorf("HLS Generation %s terminó con error: %w", res.Label, err)
	}

	return nil
}
