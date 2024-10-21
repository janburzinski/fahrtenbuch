package main

import (
	"server/pkg/logger"
)

func main() {
	err := logger.CreateLoggerDir()
	if err != nil {
		panic(err)
	}
	logger.Log("DEBUG", "Successfully initialized the Logger")
}
