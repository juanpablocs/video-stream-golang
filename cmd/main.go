package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/juanpablocs/video-stream-golang/internal/dependencies"
	"github.com/juanpablocs/video-stream-golang/internal/handlers"
	"github.com/juanpablocs/video-stream-golang/internal/routes"
	"github.com/juanpablocs/video-stream-golang/internal/usecases"
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

	usecase := usecases.NewUsecase(mongoClient.Database(os.Getenv("MONGODB_DATABASE")))
	handler := handlers.NewHandler(
		// TODO: Change to repository or usecase
		mongoClient.Database(os.Getenv("MONGODB_DATABASE")),
		channel,
		usecase,
	)
	// Create a new Fiber instance.
	app := fiber.New()
	app.Use(
		logger.New(),
	)
	routes.AddRoutes(app, handler)
	log.Fatal(app.Listen(":3001"))
}
