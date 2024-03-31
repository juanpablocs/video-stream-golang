package main

import (
	"fmt"
	"log"
	"os"

	"github.com/juanpablocs/video-stream-golang/internal/dependencies"
	"github.com/juanpablocs/video-stream-golang/internal/usecases"
	pkgVideo "github.com/juanpablocs/video-stream-golang/pkg/video"
)

func main() {
	dependencies.EnvLoad()
	// Crea una nueva conexión RabbitMQ.
	connection := dependencies.ConnectToRabbitMQ()
	defer connection.Close()

	// Abre un nuevo canal.
	channel := dependencies.OpenChannel(connection)
	defer channel.Close()

	// Subscribing to QueueService1 for getting messages.
	messages := dependencies.SetupConsumer(channel, "QueueService1")

	// Create a new MongoDB connection.
	mongoClient := dependencies.ConnectDB()

	log.Println("Waiting for messages")

	usecase := usecases.NewUsecase(mongoClient.Database(os.Getenv("MONGODB_DATABASE")))
	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
			ID := string(message.Body)

			video, err := usecase.VideoId(ID)
			if err != nil {
				log.Printf("Error obteniendo video: %v\n", err)
				return
			}
			videoUpload, err := usecase.VideoUploadId(ID)
			if err != nil {
				log.Printf("Error obteniendo video: %v\n", err)
				return
			}
			fmt.Printf("Video DB: %+v\n", video)

			videoUrl := videoUpload.Url
			outputPath := fmt.Sprintf("videos/%s", ID)
			metadata := &pkgVideo.Metadata{Video: pkgVideo.Video{Filename: ID}}
			videoInfo, err := pkgVideo.VideoInfo(videoUrl, metadata)
			if err != nil {
				fmt.Println("Error obteniendo info del video:", err)
				return
			}
			fmt.Printf("Video info: %+v\n", videoInfo)

			usecase.VideoUploadMetadataUpdated(videoUpload, videoInfo)

			usecase.CreateThumbnails(videoUrl, ID, outputPath, videoInfo.Video.Duration)

			if err := usecases.CreateSpritesAndVTT(videoUrl, outputPath, videoInfo.Video.Duration); err != nil {
				fmt.Printf("Error en el proceso: %s\n", err)
			}

			if err := usecase.TranscodeVideo(videoUrl, outputPath, videoInfo.Video.Width, videoInfo.Video.Height); err != nil {
				fmt.Printf("\nError durante los procesos: %v\n", err)
			}
		}
	}()

	<-forever
}
