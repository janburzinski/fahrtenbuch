package main

import (
	"os"
	"server/pkg/logger"
	"server/pkg/router"

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
	//err = db.Connect()
	//if err != nil {
	//		panic(err)
	//	}
	logger.Log("INFO", "Successfully connected to the MySQL Database!")

	// init routes and start web server
	router := router.InitializeRoutes()
	appPort := getAppPort()
	router.Run(appPort)
}

func getAppPort() string {
	if appPort, has := os.LookupEnv("SERVER_PORT"); !has {
		logger.Log(logger.LOG_WARN, "Environmental Variable 'SERVER_PORT' was not set, using :8080")
		return ":8080"
	} else {
		return ":" + appPort
	}
}
