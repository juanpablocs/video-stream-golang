package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/juanpablocs/ffmpeg-golang/internal/dependencies"
	"github.com/juanpablocs/ffmpeg-golang/internal/handlers"
	"github.com/juanpablocs/ffmpeg-golang/internal/routes"
)

func main() {
	fmt.Println("Starting the application...")
	dependencies.EnvLoad()
	// Crea una nueva conexión RabbitMQ.
	connection := dependencies.ConnectToRabbitMQ()
	defer connection.Close()

	// Abre un nuevo canal.
	channel := dependencies.OpenChannel(connection)
	defer channel.Close()

	dependencies.SetupQueue(channel, "QueueService1")

	// Create a new MongoDB connection.
	mongoClient := dependencies.ConnectDB()

	handler := handlers.NewHandler(
		// TODO: Change to repository or usecase
		mongoClient.Database(os.Getenv("MONGODB_DATABASE")),
		channel,
	)
	// Create a new Fiber instance.
	app := fiber.New()
	app.Use(
		logger.New(),
	)
	routes.AddRoutes(app, handler)
	log.Fatal(app.Listen(":3001"))
}
