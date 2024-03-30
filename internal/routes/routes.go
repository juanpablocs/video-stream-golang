package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/juanpablocs/ffmpeg-golang/internal/handlers"
)

func AddRoutes(app *fiber.App, handler *handlers.Handler) {
	app.Static("/bucket", "./videos")
	app.Get("/videos", handler.ListVideos)
	app.Get("/videos/:id", handler.VideoByID)
	app.Post("/videos", handler.CreateVideo)
}
