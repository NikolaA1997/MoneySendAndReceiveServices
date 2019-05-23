package main

import (
	"app/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()

	config.AllowAllOrigins = true

	config.AllowHeaders = []string{"Content-Type"}
	config.AllowMethods = []string{"POST", "GET"}
	router.Use(cors.New(config))

	router.POST("send-money", handleAndSendMoney)
	router.GET("get-money", retreiveMoney)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(400, "Page not found")
	})

	err := router.Run(":8000")
	if err != nil {
		util.ErrorLog(err)
	}
}
