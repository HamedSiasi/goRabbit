package main

import (
	"bufio"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

// Helper function to check the return value for each call
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func sendMsg(myname string, destination string, message string) bool {
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
		destination, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
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
			Body:        []byte(myname + ":" + message),
		})
	failOnError(err, "Failed to publish a message")

	return retval
}

func receiveMsg(myname string) {
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
		myname, // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
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
			log.Printf("---> %s", d.Body)
		}
	}()
	forever := make(chan bool)
	<-forever
}

func main() {
	if len(os.Args) == 2 {
		myname := os.Args[1]

		//Receiving goroutine
		go func() {
			receiveMsg(myname)
		}()

		//Publishing goroutine
		go func() {
			for {
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')

				index := strings.Index(text, ":")
				if index > 0 {
					destination := text[:index]
					message := text[index+1:]
					if len(destination) > 0 && len(message) > 0 {
						//fmt.Println("index:", index)
						//fmt.Println("destination:", destination)
						//fmt.Println("message:", message)
						sendMsg(myname, destination, message)
					} else {
						fmt.Println("Destination: Message")
					}
				} else {
					fmt.Println("Destination: Message")
				}
			}
		}()

		// keep waiting the main goroutine
		// We'll keep it running to listen for messages and print them out.
		forever := make(chan bool)
		<-forever
	} else {
		fmt.Println("Client name is missing!")
	}
} //main
