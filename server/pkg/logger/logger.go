package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LOG_INFO  = "INFO"
	LOG_DEBUG = "DEBUG"
	LOG_ERROR = "ERROR"
)

var (
	filePath = "logs/log.json"
)

type LoggerMessage struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// automatically log to json file whenever a panic occurs
func RecoverAndLog() {

}

func Log(level string, message string) {
	//write log to file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//validate logger level
	if level != LOG_INFO && level != LOG_DEBUG && level != LOG_ERROR {
		log.Printf("ERROR: Got a wrong LOG LEVEL. Got: %s, Changed to: %s", level, "INFO")
		level = "INFO"
	}

	//build json struct
	logMsg := LoggerMessage{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	jsonMessage, err := json.Marshal(logMsg)
	if err != nil {
		panic(err)
	}

	// Read the existing content of the log file
	content, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	var newContent string
	if len(content) == 0 {
		//if the file is empty, start a new JSON array
		newContent = "[\n" + string(jsonMessage) + "\n]"
	} else {
		//if the file is not empty, append the new log message
		contentStr := strings.TrimSpace(string(content))
		if strings.HasSuffix(contentStr, "]") {
			//remove the closing bracket
			contentStr = contentStr[:len(contentStr)-1]
			newContent = contentStr + ",\n" + string(jsonMessage) + "\n]"
		} else {
			panic("Invalid log file format")
		}
	}

	// write the updated content back to the file
	// but dont log debug messages if in prod env
	currEnv := os.Getenv("GO_ENV")
	if currEnv != "debug" {
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			panic(err)
		}

		//also output into console
		fmt.Println(string(jsonMessage))
	}
}

func CreateLoggerDir() error {
	err := os.MkdirAll("logs", os.ModePerm)
	return err
}
