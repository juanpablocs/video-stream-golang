package usecases

import (
	"context"

	"github.com/juanpablocs/ffmpeg-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u Usecase) VideoStatus(ID string, status models.VideoStatus) error {
	collection := u.db.Collection("videos")
	idObj, _ := primitive.ObjectIDFromHex(ID)
	updateVideo := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": idObj}, updateVideo)
	if err != nil {
		return err
	}
	return nil
}
