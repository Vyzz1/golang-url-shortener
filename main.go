package main

import (
	"log"
	"url-shortener/utils"

	"github.com/gin-gonic/gin"
)

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "OK",
		})
	})

	log.Printf("Starting HTTP server on %s", config.HttpServerAddress)
	router.Run(config.HttpServerAddress)
}
