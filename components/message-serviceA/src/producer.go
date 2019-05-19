package main

import (
	"encoding/json"
	"fmt"
	"github.com/MoneySendAndReceiveServices/components/message-serviceA/src/models"
	"github.com/MoneySendAndReceiveServices/components/message-serviceA/src/util"
	"github.com/streadway/amqp"
)

func produceMessage(message models.MessagePayload) {

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

	body, _ := json.Marshal(message)

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	util.FailOnError(err, "Failed to publish a message")
	fmt.Println("Message sucessfuly produced")

	return
}
