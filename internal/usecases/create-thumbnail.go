package usecases

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/juanpablocs/ffmpeg-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Usecase) CreateThumbnails(videoUrl, videoID, videoDir string, duration float64) {
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		fmt.Printf("creating directory: %v", err)
		return
	}

	times := []float64{duration * 0.1, duration * 0.5, duration * 0.9} // Calcular 10%, 50%, 90%
	positions := []string{"start", "middle", "end"}
	var thumbnails []models.Thumbnail // Inicializa un slice vacío para almacenar los thumbnails

	for pos, t := range times {
		timestamp := fmt.Sprintf("%f", t)
		smallPath := fmt.Sprintf("%s/small-%d.jpg", videoDir, int(t))
		largePath := fmt.Sprintf("%s/large-%d.jpg", videoDir, int(t))

		// Ejecutar ffmpeg una vez por cada tamaño de thumbnail requerido
		cmdSmall := exec.Command("ffmpeg", "-y", "-ss", timestamp, "-i", videoUrl, "-vframes", "1", "-q:v", "2", "-vf", "scale=320:-1", smallPath)
		if err := cmdSmall.Run(); err != nil {
			fmt.Println("Error al generar thumbnail small:", err)
			return
		}
		cmdLarge := exec.Command("ffmpeg", "-y", "-ss", timestamp, "-i", videoUrl, "-vframes", "1", "-q:v", "2", "-vf", "scale=720:-1", largePath)
		if err := cmdLarge.Run(); err != nil {
			fmt.Println("Error al generar thumbnail large:", err)
			return
		}

		thumbnailSmall := models.Thumbnail{
			Size:     "small",
			Position: positions[pos],
			Path:     strings.Replace(smallPath, "videos/", "/bucket/", 1),
			Default:  pos == 1,
		}
		thumbnailLarge := models.Thumbnail{
			Size:     "large",
			Position: positions[pos],
			Path:     strings.Replace(largePath, "videos/", "/bucket/", 1),
			Default:  pos == 1,
		}

		thumbnails = append(thumbnails, thumbnailSmall, thumbnailLarge)

		fmt.Printf("Thumbnails generados para el tiempo %s: %s, %s\n", timestamp, smallPath, largePath)
	}

	// Actualizar la colección con el slice de thumbnails después de salir del bucle.
	objID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		fmt.Println("Error al convertir videoID a ObjectID:", err)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"thumbnails": thumbnails,
			"status":     models.StatusConverting,
			"updatedAt":  time.Now(),
		},
	}
	collection := h.db.Collection("videos")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		fmt.Println("Error al actualizar los thumbnails en la base de datos:", err)
	}
}

// func CreateThumbnails(videoPath, baseName, timestamp string) {
// 	thumbnailPathSmall := fmt.Sprintf("videos-m3u8/%s/thumbnail_small.jpg", baseName)
// 	thumbnailPathLarge := fmt.Sprintf("videos-m3u8/%s/thumbnail_large.jpg", baseName)

// 	cmd := exec.Command("ffmpeg", "-ss", timestamp, "-i", videoPath,
// 		"-filter_complex", "[0:v]scale=320:-1[small];[0:v]scale=720:-1[large]",
// 		"-map", "[small]", "-vframes", "1", "-q:v", "2", thumbnailPathSmall,
// 		"-map", "[large]", "-vframes", "1", "-q:v", "2", thumbnailPathLarge)

// 	if err := cmd.Run(); err != nil {
// 		fmt.Println("Error al crear thumbnails:", err)
// 		return
// 	}

// 	fmt.Println("Thumbnails creados con éxito:", thumbnailPathSmall, "y", thumbnailPathLarge)
// }
