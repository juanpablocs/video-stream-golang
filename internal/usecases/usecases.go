package usecases

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Usecase struct {
	db *mongo.Database
}

func NewUsecase(db *mongo.Database) *Usecase {
	return &Usecase{
		db: db,
	}
}
