package usecases

import (
	"context"
	"time"

	"github.com/juanpablocs/video-stream-golang/internal/models"
	"github.com/juanpablocs/video-stream-golang/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u Usecase) VideoList(page, pageSize int64) (*models.VideoListResponse, error) {
	var collection = u.db.Collection("videos")
	var videos []models.Video

	// Opciones para ordenamiento y paginación
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}). // Ordenamiento descendente por fecha de creación
		SetSkip((page - 1) * pageSize).                 // Saltar los documentos de páginas anteriores
		SetLimit(pageSize)                              // Limitar el número de documentos

	totalVideos, err := collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	filter := bson.D{{}}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &videos); err != nil {
		return nil, err
	}

	var publicVideos []models.PublicVideo
	for _, video := range videos {
		// Convertir el tiempo a string
		createdAt := video.CreatedAt.Format(time.RFC3339)

		// Construir el mapa de thumbnails
		thumbnailMap := utils.ThumbnailMap(video.Thumbnails)

		// Agregar el video convertido a la lista de PublicVideo
		publicVideos = append(publicVideos, models.PublicVideo{
			ID:          video.ID.Hex(),
			Title:       video.Title,
			Description: video.Description,
			Thumbnail:   thumbnailMap, // Ajustado para usar el mapa
			Duration:    video.Duration,
			VttTrack:    video.VttTrack,
			Status:      video.Status,
			CreatedAt:   createdAt,
		})
	}

	// Calcular la página siguiente, si hay más elementos
	nextPage := int64(0)
	if page*pageSize < totalVideos {
		nextPage = page + 1
	}

	return &models.VideoListResponse{
		CurrentPage: page,
		NextPage:    nextPage,
		Total:       totalVideos,
		Videos:      publicVideos,
	}, nil
}
