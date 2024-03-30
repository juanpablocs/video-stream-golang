package dependencies

import (
	"os"

	"github.com/streadway/amqp"
)

// ConnectToRabbitMQ establece una conexión con RabbitMQ.
func ConnectToRabbitMQ() *amqp.Connection {
	connectRabbitMQ, err := amqp.Dial(os.Getenv("AMQP_SERVER_URL"))
	if err != nil {
		panic(err)
	}
	return connectRabbitMQ
}

// OpenChannel abre un canal sobre la conexión dada de RabbitMQ.
func OpenChannel(connectRabbitMQ *amqp.Connection) *amqp.Channel {
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	return channelRabbitMQ
}

// SetupQueue declara una cola en el canal dado y retorna la cola.
func SetupQueue(channel *amqp.Channel, queueName string) amqp.Queue {
	queue, err := channel.QueueDeclare(
		queueName, // nombre de la cola
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // argumentos
	)
	if err != nil {
		panic("Failed to declare a queue: " + err.Error())
	}
	return queue
}

// SetupConsumer configura un consumidor en el canal y cola dados.
func SetupConsumer(channel *amqp.Channel, queueName string) <-chan amqp.Delivery {
	msgs, err := channel.Consume(
		queueName, // cola
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		panic("Failed to register a consumer: " + err.Error())
	}
	return msgs
}

// Publish envía un mensaje a una cola en RabbitMQ.
func Publish(channel *amqp.Channel, messageBody string) error {
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(messageBody),
	}

	// Intentar publicar un mensaje en la cola.
	err := channel.Publish(
		"",              // exchange
		"QueueService1", // nombre de la cola
		false,           // mandatory
		false,           // immediate
		message,         // mensaje a publicar
	)
	return err
}
