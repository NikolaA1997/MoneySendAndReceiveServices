package main

import (
	"encoding/json"
	"flag"
	"github.com/MoneySendAndReceiveServices/components/message-serviceB/src/models"
	"github.com/MoneySendAndReceiveServices/components/message-serviceB/src/util"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {

	storageFile := flag.String("storageFile", "../storage.json", "JSON storage file")
	storage, err := InitStorage(storageFile)
	if err != nil {
		util.ErrorLog(err)
		return
	}

	go setupRoutes(storage)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"money", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	util.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	util.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := storeMessage(d.Body, storage)
			if err != nil {
				util.ErrorLog(err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}


func storeMessage(body []byte, storage *Storage) error {
	message := models.Message{}

	if err := json.Unmarshal(body, &message); err != nil {
		return err
	}

	if err := storage.Set(message.Amount, time.Now()); err != nil {
		return err
	}

	return nil
}

