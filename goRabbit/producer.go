package main

import (
	"bufio"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

// Helper function to check the return value for each call
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func sendMsg(body string) bool {
	retval := true

	// (1) --- Open a Connection ---
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// (2) --- Open a channel ---
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// (3) --- Declare a queue ---
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// (4) ---- Publish on queue ----
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	return retval
}

func main() {
	//var str string

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		//fmt.Println(text)
		result := sendMsg(text)
		if !result {
			fmt.Println("error")
		} else {
			fmt.Println("done")
		}
	} //for

} //main

// Broker server
// go https://www.rabbitmq.com/install-windows.html
// Hamed Siasi
