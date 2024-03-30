package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/juanpablocs/ffmpeg-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h Handler) ListVideos(c *fiber.Ctx) error {
	// Crea un contexto con un timeout.
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel() // Asegura que el contexto se cancele para liberar recursos.

	collection := h.db.Collection("videos")
	var videos []models.Video

	filter := bson.D{{}}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err = cursor.All(ctx, &videos); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(videos)
}
