package main

import (
	"app/models"
	"app/util"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}

type stop struct {
	error
}

func main() {

	storageFile := flag.String("storageFile", "../storage.json", "JSON storage file")
	storage, err := InitStorage(storageFile)
	if err != nil {
		util.ErrorLog(err)
		return
	}

	go setupRoutes(storage)

	duration, _ := time.ParseDuration("10s")
	retry(3, duration, func() (err error) {
		conn, err := amqp.Dial("amqp://guest:guest@rabbit:5672/")
		util.FailOnError(err, "Failed to connect to RabbitMQ")

		if err == nil {

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

			fmt.Println("Channel and Queue established")
			fmt.Println(q)

			defer conn.Close()
			defer ch.Close()

			msgs, err := ch.Consume(
				q.Name, // queue
				"",     // consumer
				false,   // auto-ack
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
					d.Ack(false)
				}
			}()

			log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
			<-forever
		}
		return err
	})
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
