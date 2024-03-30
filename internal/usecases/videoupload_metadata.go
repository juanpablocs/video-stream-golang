package usecases

import (
	"context"
	"time"

	"github.com/juanpablocs/video-stream-golang/internal/models"
	pkgVideo "github.com/juanpablocs/video-stream-golang/pkg/video"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Usecase) VideoUploadMetadataUpdated(v *models.VideoUpload, m *pkgVideo.Metadata) error {
	updateVideo := bson.M{
		"$set": bson.M{
			"duration":  m.Video.Duration,
			"updatedAt": time.Now(), // Actualizar la fecha de modificación
		},
	}
	collection := h.db.Collection("videos")
	collection.UpdateOne(context.TODO(), bson.M{"_id": v.VideoID}, updateVideo)

	updateVideoUpload := bson.M{
		"$set": bson.M{
			"filesize":         m.Video.FileSize,
			"videoCodec":       m.Video.Codec,
			"videoBitrate":     m.Video.BitRate,
			"videoWidth":       m.Video.Width,
			"videoHeight":      m.Video.Height,
			"videoAspectRatio": m.Video.AspectRatio,
			"audioCodec":       m.Audio.Codec,
			"audioChannels":    m.Audio.Channels,
			"audioSampleRate":  m.Audio.SampleRate,
		},
	}
	collectionUpload := h.db.Collection("videoUploads")
	collectionUpload.UpdateOne(context.TODO(), bson.M{"_id": v.ID}, updateVideoUpload)
	return nil
}
