package main

import (

	"fmt"
	"app/models"
	"app/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func handleAndSendMoney(c *gin.Context) {

	var message models.MessagePayload

	bindErr := c.BindJSON(&message)
	if bindErr != nil {
		util.ErrorLog(bindErr)
		c.JSON(400, "Invalid Data")
		return
	}

	if validateErr := message.Validate(); validateErr != nil {
		util.ErrorLog(validateErr)
		c.JSON(400, "Invalid Data")
		return
	}
	fmt.Println(message)
	message.ConvertAmount()
	produceMessage(message)

}

func retreiveMoney(c *gin.Context) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://message-serviceb:8001/storage", nil)
	resp, err := client.Do(req)
	if err != nil {
		util.ErrorLog(err)
		c.JSON(500, "Internal error")
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	c.JSON(200, string(body))
	return
}

