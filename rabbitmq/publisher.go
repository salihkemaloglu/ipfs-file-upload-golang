package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func PublishDataToQueue(message string) (string,error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!=nil{
		return "Failed to connect to RabbitMQ",err
	}
	//failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	if err!=nil{
		return "Failed to open a channel",err
	}
	//failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		true,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	//failOnError(err, "Failed to declare a queue")
	if err!=nil{
		return "Failed to declare a queue",err
	}
	body := message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	//log.Printf(" [x] Sent %s", body)
	if err !=nil{
		return "Failed to publish a message",err
	}
	return "File addet to ipfs sucsessfully.File hash code is",nil
	//failOnError(err, "Failed to publish a message")
}