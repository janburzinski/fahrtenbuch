package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var (
	filePath = "logs/log.json"
)

type LoggerMessage struct {
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func Log(level string, message string) {
	//write log to file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

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
		// If the file is empty, start a new JSON array
		newContent = "[\n" + string(jsonMessage) + "\n]"
	} else {
		// If the file is not empty, append the new log message
		contentStr := strings.TrimSpace(string(content))
		if strings.HasSuffix(contentStr, "]") {
			// Remove the closing bracket
			contentStr = contentStr[:len(contentStr)-1]
			newContent = contentStr + ",\n" + string(jsonMessage) + "\n]"
		} else {
			// Handle unexpected content
			panic("Invalid log file format")
		}
	}

	// Write the updated content back to the file
	err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		panic(err)
	}

	//also output into console
	fmt.Println(string(jsonMessage))
}

func CreateLoggerDir() error {
	err := os.MkdirAll("logs", os.ModePerm)
	return err
}
