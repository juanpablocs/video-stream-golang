package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/juanpablocs/video-stream-golang/internal/dependencies"
	"github.com/juanpablocs/video-stream-golang/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h Handler) CreateVideo(c *fiber.Ctx) error {
	var videoRequest models.VideoRequest
	if err := c.BodyParser(&videoRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error al parsear el body",
			"error":   err,
		})
	}

	// Insertar el video en MongoDB.
	collection := h.db.Collection("videos")
	newVideo := models.Video{
		ID:          primitive.NewObjectID(),
		Title:       videoRequest.Title,
		Description: videoRequest.Description,
		Status:      models.StatusIdle,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := collection.InsertOne(context.TODO(), newVideo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al insertar el video en la base de datos",
			"error":   err,
		})
	}

	videoUpload := models.VideoUpload{
		VideoID: newVideo.ID,
		Url:     videoRequest.Url,
	}

	uploadCollection := h.db.Collection("videoUploads")
	_, err = uploadCollection.InsertOne(context.TODO(), videoUpload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al insertar el videoUpload en la base de datos",
			"error":   err,
		})
	}

	if err := dependencies.Publish(h.channel, newVideo.ID.Hex()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error al publicar el video en la cola de mensajes",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Video creado exitosamente",
		"video":   newVideo,
	})

}
