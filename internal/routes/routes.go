package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/juanpablocs/video-stream-golang/internal/handlers"
)

func AddRoutes(app *fiber.App, handler *handlers.Handler) {
	app.Static("/bucket", "./videos")
	app.Get("/videos", handler.ListVideos)
	app.Get("/videos/:id", handler.VideoByID)
	app.Post("/videos", handler.CreateVideo)

	// Configurar CORS para rutas públicas
	corsConfig := cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST",
	})
	public := app.Group("/v")
	public.Use(corsConfig)
	public.Get("/", handler.PublicListVideos)
	public.Get("/:id", handler.PublicVideoID)
}
