package main

import (
	"github.com/MoneySendAndReceiveServices/components/message-serviceB/src/util"
	"github.com/gin-gonic/gin"
)

func setupRoutes(storage *Storage) {
	router := gin.Default()

	router.GET("storage", func(c *gin.Context) {
		money := storage.Get()
		c.JSON(200, money)
		return
	})
	err := router.Run(":8001")
	if err != nil {
		util.ErrorLog(err)
		return
	}

}
