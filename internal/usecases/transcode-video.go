package usecases

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/juanpablocs/video-stream-golang/internal/models"
	"github.com/juanpablocs/video-stream-golang/internal/utils"
	pkgVideo "github.com/juanpablocs/video-stream-golang/pkg/video"
)

func (u Usecase) TranscodeVideo(originalInputPath, outputPath string, width, height int) error {
	currentInputPath := originalInputPath // Empieza con el path original

	masterPlaylistPath := filepath.Join(outputPath, "video_master.m3u8")
	var variantPlaylistEntries []string

	nameSplits := strings.Split(outputPath, "/")
	ID := nameSplits[len(nameSplits)-1]

	applicableResolutions := utils.FilterResolutions(width, height)

	for i, res := range applicableResolutions {
		resOutputMP4 := fmt.Sprintf("%s/video_%s.mp4", outputPath, res.Label)
		resOutputHLS := fmt.Sprintf("%s/video_%s.m3u8", outputPath, res.Label)

		// Crear las carpetas necesarias
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			u.VideoStatus(ID, models.StatusError) //TODO: add log error to database
			return fmt.Errorf("error al crear las carpetas necesarias: %w", err)
		}

		// Ejecutar la transcodificación para la resolución actual
		if err := executeAndMonitor(currentInputPath, resOutputMP4, resOutputHLS, res); err != nil {
			u.VideoStatus(ID, models.StatusError) //TODO: add log error to database
			return err
		}

		// Para la siguiente iteración, usar el MP4 que acabamos de crear como input
		currentInputPath = resOutputMP4

		// Para evitar re-codificar el original en la primera iteración, restaurar el inputPath original
		if i == 0 {
			currentInputPath = originalInputPath
		}
		// Agregar entrada de playlist de variante al archivo maestro
		bandwidth := calculateBandwidth(res)
		variantPlaylistEntries = append(variantPlaylistEntries, fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,CODECS=\"avc1.42e01e,mp4a.40.2\",RESOLUTION=%dx%d\n%s", bandwidth, res.Width, res.Height, fmt.Sprintf("video_%s.m3u8", res.Label)))

		// get video info metadata
		metadata := &pkgVideo.Metadata{Video: pkgVideo.Video{Filename: ID}}
		videoInfo, errVideoInfo := pkgVideo.VideoInfo(resOutputMP4, metadata)
		if errVideoInfo != nil {
			u.VideoStatus(ID, models.StatusError) //TODO: add log error to database
			return fmt.Errorf("error obteniendo info del video: %w", errVideoInfo)
		}
		// insert the transcoded video into the database
		if errInsert := u.InsertVideoTranscoded(ID, resOutputMP4, resOutputHLS, videoInfo); errInsert != nil {
			u.VideoStatus(ID, models.StatusError) //TODO: add log error to database
			return fmt.Errorf("error al insertar el video transcodificado: %w", errInsert)
		}
	}

	// Crear el archivo maestro M3U8 que incluye todas las variantes
	masterContent := "#EXTM3U\n" + strings.Join(variantPlaylistEntries, "\n")
	if err := os.WriteFile(masterPlaylistPath, []byte(masterContent), 0644); err != nil {
		u.VideoStatus(ID, models.StatusError) //TODO: add log error to database
		return fmt.Errorf("error al escribir el archivo maestro M3U8: %w", err)
	}

	u.VideoStatus(ID, models.StatusFinished)
	fmt.Println("\nProcesos completados con éxito.")
	return nil
}

func executeAndMonitor(inputPath, outputPathMP4, outputPathHLS string, res models.Resolution) error {
	var wg sync.WaitGroup
	wg.Add(1)

	// Generar MP4
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

func calculateBandwidth(res models.Resolution) int {
	// ancho de banda aproximado basado en la resolución y otros factores
	return res.Width * res.Height * 5
}
