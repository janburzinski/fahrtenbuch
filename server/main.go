package main

import (
	"server/pkg/db"
	"server/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//load dotenv
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	//init the logger
	err = logger.CreateLoggerDir()
	if err != nil {
		panic(err)
	}
	logger.Log("DEBUG", "Successfully initialized the Logger")

	//init db
	err = db.Connect()
	if err != nil {
		panic(err)
	}
	logger.Log("INFO", "Successfully connected to the MySQL Database!")

	// init web server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() //:8080
}
