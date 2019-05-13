package main

import (
	"github.com/MoneySendAndReceiveServices/storage"
	"github.com/MoneySendAndReceiveServices/util"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)


func main() {
	TOKENN := "5"
	store := storage.NewStorage(".storage.json")
	err := store.Set(0, time.Time{})


		router := gin.Default()

		router.GET("storage", func(c *gin.Context) {
			token := c.GetHeader("Authorization")
			if token != TOKENN {
				c.JSON(400, "UNAUTHORIZED")
				return
			}
				money, err := store.Get()
				if err != nil {
					util.ErrorLog(err)
				}
				c.JSON(200, money)
				return
			})
		err = router.Run(":8001")
		if err != nil {
			util.ErrorLog(err)
		}


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
			messageWords := util.MessageFormat(string(d.Body))
			amount, err := strconv.ParseFloat(messageWords[0],64)
			if err != nil {
				util.ErrorLog(err)
			}
			err = store.Set(amount,time.Now())
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
