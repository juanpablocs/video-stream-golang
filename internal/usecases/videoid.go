package usecases

import (
	"context"
	"fmt"

	"github.com/juanpablocs/ffmpeg-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h Usecase) VideoId(id string) (*models.Video, error) {
	video := &models.Video{}

	// Convertir el string id a un ObjectID
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Error al convertir id a ObjectID:", err)
		return nil, err // Es importante manejar este error
	}
	// Asegúrate de usar objId en tu filtro y el campo correcto "_id"
	v := h.db.Collection("videos").FindOne(context.TODO(), bson.M{"_id": objId})
	if v.Err() != nil {
		return nil, v.Err()
	}

	// Decodificar el resultado en la estructura video
	if err := v.Decode(video); err != nil {
		return nil, err
	}

	return video, nil
}
