package logger

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	// Initialize logger
	err := CreateLoggerDir()
	if err != nil {
		t.Fatalf("Error while creating logger dir: %s", err.Error())
	}

	tmpFile, err := os.CreateTemp("logs", "testlog")
	if err != nil {
		t.Fatalf("Error while creating temporary file: %s", err.Error())
	}
	defer os.Remove(tmpFile.Name())

	// Override the already existing file path to the tmp file for testing purposes
	filePath = tmpFile.Name()

	// Now actually test the logger
	testLogMessage := "test"
	Log("INFO", testLogMessage)

	// Open log file
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error while reading test file: %s", err.Error())
	}

	// Check the content of the log message
	var logMessages []LoggerMessage
	err = json.Unmarshal(content, &logMessages)
	if err != nil {
		t.Fatalf("Error while converting file content to json: %s", err.Error())
	}

	if len(logMessages) == 0 || logMessages[0].Level != "INFO" || logMessages[0].Message != testLogMessage {
		t.Fatalf("The content of the log file were wrong. expected: %s, %s, got: %s, %s", "INFO", testLogMessage, logMessages[0].Level, logMessages[0].Message)
	}
}
