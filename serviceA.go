package main

import (
	"./util"
	"fmt"
	"github.com/MoneySendAndReceiveServices/message"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
)


func main() {
	router := gin.Default()

	config := cors.DefaultConfig()

	config.AllowAllOrigins = true

	config.AllowHeaders = []string{"Content-Type"}
	config.AllowMethods = []string{"POST", "GET"}
	router.Use(cors.New(config))

	router.POST("send-money", handleMoney)
	router.GET("get-money", retreiveMoney)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(400, "Page not found")
	})

	err := router.Run(":8000")
	if err != nil {
		util.ErrorLog(err)
	}
}

func handleMoney(c *gin.Context) {

	var mess message.MessagePayload

	bindErr := c.BindJSON(&mess)
	if bindErr != nil {
		util.ErrorLog(bindErr)
		c.JSON(400, "Invalid Data")
		return
	}

	if validateErr := mess.Validate(); validateErr != nil {
		util.ErrorLog(validateErr)
		c.JSON(400, "Invalid Data")
		return
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

	body := mess
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(fmt.Sprintf("%v", mess)),
		})
	log.Printf(" [x] Sent %s", body)
	util.FailOnError(err, "Failed to publish a message")

	fmt.Println(mess)
	return
}

func retreiveMoney(c *gin.Context){

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8001/storage", nil)
	req.Header.Add("Authorization", "5")
	resp, err := client.Do(req)
	if resp.StatusCode != 200{
		c.JSON(400, "Sorry brah")
		return
	}
	if err != nil {
		util.ErrorLog(err)
		c.JSON(500, "Internal error")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	c.JSON(200, fmt.Sprintf("%s", body))
	return
}
