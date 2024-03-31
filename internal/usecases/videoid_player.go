package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/juanpablocs/video-stream-golang/internal/models"
	"github.com/juanpablocs/video-stream-golang/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h Usecase) VideoIdPlayer(id string) (*models.PublicVideo, error) {
	video := &models.Video{}
	videoPublic := &models.PublicVideo{}
	videosTranscoded := []models.VideoUploadTranscoded{}
	// Convertir el string id a un ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error al convertir id a ObjectID:", err)
		return nil, err
	}

	v := h.db.Collection("videos").FindOne(context.TODO(), bson.M{"_id": objId})
	if v.Err() != nil {
		return nil, v.Err()
	}
	if err := v.Decode(video); err != nil {
		return nil, err
	}

	cursor, errTrans := h.db.Collection("videoUploadsTranscoded").Find(context.TODO(), bson.M{"videoId": objId})
	if errTrans != nil {
		return nil, errTrans
	}
	if err = cursor.All(context.TODO(), &videosTranscoded); err != nil {
		return nil, err
	}

	sources := []string{}
	for _, videoTranscoded := range videosTranscoded {
		sources = append(sources, strings.Replace(videoTranscoded.UrlMP4, "videos/", "/bucket/", 1))
	}

	videoPublic.ID = video.ID.Hex()
	videoPublic.Title = video.Title
	videoPublic.Description = video.Description
	videoPublic.Duration = video.Duration
	videoPublic.Thumbnail = utils.ThumbnailMap(video.Thumbnails)
	videoPublic.VttTrack = fmt.Sprintf("/bucket/%s/sprites/thumbnails.vtt", videoPublic.ID)
	videoPublic.Status = video.Status
	videoPublic.CreatedAt = video.CreatedAt.Format("2006-01-02T15:04:05")
	videoPublic.Sources = sources
	videoPublic.Playlist = fmt.Sprintf("/bucket/%s/video_master.m3u8", videoPublic.ID)
	return videoPublic, nil
}
