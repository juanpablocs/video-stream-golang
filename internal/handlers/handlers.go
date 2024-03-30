package handlers

import (
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	db      *mongo.Database
	channel *amqp.Channel
}

func NewHandler(db *mongo.Database, channel *amqp.Channel) *Handler {
	return &Handler{
		db:      db,
		channel: channel,
	}
}
