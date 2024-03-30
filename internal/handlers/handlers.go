package handlers

import (
	"github.com/juanpablocs/ffmpeg-golang/internal/usecases"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	db      *mongo.Database
	channel *amqp.Channel
	usecase *usecases.Usecase
}

func NewHandler(db *mongo.Database, channel *amqp.Channel, usecase *usecases.Usecase) *Handler {
	return &Handler{
		db:      db,
		channel: channel,
		usecase: usecase,
	}
}
