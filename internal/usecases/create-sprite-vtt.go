package usecases

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/juanpablocs/video-stream-golang/pkg/util"
)

const (
	cols       = 4   // Número de columnas en el sprite
	rows       = 8   // Número de filas en el sprite
	thumbWidth = 200 // Ancho fijo de los thumbnails
)

// TODO: refactor algorithm for calculating the interval and video duration
func CreateSpritesAndVTT(videoPath string, outputDir string, videoDuration float64) error {
	spritesDir := filepath.Join(outputDir, "sprites")
	secondInterval := util.CalculeInterval(videoDuration / 60)
	if err := os.MkdirAll(spritesDir, 0755); err != nil {
		return fmt.Errorf("creating sprites directory: %w", err)
	}

	// Generar todos los thumbnails con FFmpeg
	cmd := exec.Command("ffmpeg", "-i", videoPath,
		"-vf", fmt.Sprintf("fps=1/%d,scale='min(%d,iw):-1'", int(secondInterval), thumbWidth),
		"-q:v", "2", //Los valores típicos están entre 2 (alta calidad) y 31 (baja calidad)
		filepath.Join(spritesDir, "thumb-%04d.jpg"))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generating thumbnails with FFmpeg: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(spritesDir, "*.jpg"))
	if err != nil {
		return fmt.Errorf("listing thumbnail files: %w", err)
	}

	var thumbHeight int
	if len(files) > 0 {
		img, err := gg.LoadImage(files[0])
		if err != nil {
			return fmt.Errorf("loading first thumbnail: %w", err)
		}
		thumbHeight = img.Bounds().Dy()
	}

	spriteWidth := cols * thumbWidth
	spriteHeight := rows * thumbHeight
	thumbsPerSprite := cols * rows
	numSprites := (len(files) + thumbsPerSprite - 1) / thumbsPerSprite

	// Crear sprites y secciones VTT por cada segmento de thumbnails
	for s := 0; s < numSprites; s++ {
		start := s * thumbsPerSprite
		end := start + thumbsPerSprite
		if end > len(files) {
			end = len(files)
		}
		segmentFiles := files[start:end]

		// Crear el sprite para el segmento actual
		spritePath, err := createSprite(segmentFiles, s, spriteWidth, spriteHeight, thumbWidth, thumbHeight, spritesDir)
		if err != nil {
			return err
		}

		// Agregar entradas al archivo VTT para el sprite actual
		if err := appendToVTTFile(spritePath, segmentFiles, start, secondInterval, videoDuration, spritesDir, thumbHeight); err != nil {
			return err
		}
	}

	fmt.Println("Thumbnails, Sprites, and VTT file created successfully.")
	return nil
}

// createSprite crea un sprite con los thumbnails proporcionados y devuelve la ruta del sprite creado.
func createSprite(files []string, spriteIndex, spriteWidth, spriteHeight, thumbWidth, thumbHeight int, outputDir string) (string, error) {
	dc := gg.NewContext(spriteWidth, spriteHeight)
	for i, file := range files {
		img, err := gg.LoadImage(file)
		if err != nil {
			return "", fmt.Errorf("loading image: %w", err)
		}
		x := (i % cols) * thumbWidth
		y := (i / cols) * thumbHeight
		dc.DrawImageAnchored(img, x, y, 0, 0)
	}
	spritePath := filepath.Join(outputDir, fmt.Sprintf("sprite-%d.png", spriteIndex))
	if err := dc.SavePNG(spritePath); err != nil {
		return "", fmt.Errorf("saving sprite: %w", err)
	}
	return spritePath, nil
}

// appendToVTTFile agrega las entradas VTT para los thumbnails en un sprite específico.
func appendToVTTFile(spritePath string, files []string, startIdx int, secondInterval float64, videoDuration float64, outputDir string, thumbHeight int) error {
	vttPath := filepath.Join(outputDir, "thumbnails.vtt")
	file, err := os.OpenFile(vttPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening VTT file: %w", err)
	}
	defer file.Close()

	for i := range files {
		thumbIdx := startIdx + i
		startTime := float64(thumbIdx) * secondInterval
		endTime := startTime + secondInterval
		if endTime > videoDuration {
			endTime = videoDuration
		}

		// Convertir segundos a formato VTT MM:SS.mmm
		startVTT := convertToVTTTime(startTime)
		endVTT := convertToVTTTime(endTime)

		column := thumbIdx % cols
		row := (thumbIdx / cols) % rows
		x := column * thumbWidth
		y := row * thumbHeight

		spriteFileName := filepath.Base(spritePath)
		vttEntry := fmt.Sprintf("%s --> %s\n%s#xywh=%d,%d,%d,%d\n\n", startVTT, endVTT, spriteFileName, x, y, thumbWidth, thumbHeight)

		if _, err := file.WriteString(vttEntry); err != nil {
			return fmt.Errorf("writing to VTT file: %w", err)
		}
	}

	return nil
}

// convertToVTTTime convierte segundos a formato VTT MM:SS.mmm.
func convertToVTTTime(seconds float64) string {
	mins := int(seconds) / 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(mins*60) - float64(secs)) * 1000)
	return fmt.Sprintf("%02d:%02d.%03d", mins, secs, millis)
}
