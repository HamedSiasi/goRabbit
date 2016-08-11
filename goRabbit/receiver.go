package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// Helper function to check the return value for each call
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {

	// (1) --- Open connection ----
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// (2) ---- Open channel ----
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// (3) --- Declare a Queue ---
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// (4) Tell the server to deliver us the messages from the queue.
	// Since it will push us messages asynchronously,
	// we will read the messages from a channel
	// (returned by amqp::Consume) in a goroutine.

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	// keep waiting the main goroutine
	// We'll keep it running to listen for messages and print them out.
	forever := make(chan bool)
	<-forever
}

// Broker server
// go https://www.rabbitmq.com/install-windows.html
// Hamed Siasi
