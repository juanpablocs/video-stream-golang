package usecases

import (
	"context"

	"github.com/juanpablocs/ffmpeg-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h Usecase) VideoUploadId(videoId string) (*models.VideoUpload, error) {
	// Convertir el string id a un ObjectID
	objId, _ := primitive.ObjectIDFromHex(videoId)

	var videoUpload *models.VideoUpload
	uploadCollection := h.db.Collection("videoUploads")
	err := uploadCollection.FindOne(context.TODO(), bson.M{"videoId": objId}).Decode(&videoUpload)

	if err != nil {
		return nil, err
	}

	return videoUpload, nil
}
