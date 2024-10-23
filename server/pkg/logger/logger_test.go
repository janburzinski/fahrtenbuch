package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func setupTest(t *testing.T) {
	//"load env file"
	os.Setenv("GO_ENV", "testing")

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
}

func TestLogger(t *testing.T) {
	setupTest(t)

	// Now actually test the logger
	testLogMessage := "test"
	Log(LOG_INFO, testLogMessage)

	// Check the content of the log message
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error while reading test file: %s", err.Error())
	}

	var logMessages []LoggerMessage //this has to be an array because the of the way the json struct gets saved into the json file
	err = json.Unmarshal(content, &logMessages)
	if err != nil {
		t.Fatalf("Error while converting file content to json: %s", err.Error())
	}

	if len(logMessages) == 0 || logMessages[0].Level != "INFO" || logMessages[0].Message != testLogMessage {
		t.Fatalf("The content of the log file were wrong. expected: %s, %s, got: %s, %s", "INFO", testLogMessage, logMessages[0].Level, logMessages[0].Message)
	}
}

func TestLoggerWithInvalidLe(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	setupTest(t)

	testLogMessage := "test should fail"
	Log("miau", testLogMessage)

	//check if i got a console message
	expectedErrorMessage := fmt.Sprintf("ERROR: Got a wrong LOG LEVEL. Got: %s, Changed to: %s", "miau", "INFO")
	if !bytes.Contains(buf.Bytes(), []byte(expectedErrorMessage)) {
		t.Fatalf("Expected console error message was invalid. Got: %q but go %q", expectedErrorMessage, buf.String())
	}

	//check if the LOG_LEVEL got changes to "INFO"
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error while reading test log file: %s", err.Error())
	}

	var logMessages []LoggerMessage
	err = json.Unmarshal(content, &logMessages)
	if err != nil {
		t.Fatalf("Error while converting file content to json: %s", err.Error())
	}

	if len(logMessages) == 0 || logMessages[0].Level != "INFO" || logMessages[0].Message != testLogMessage {
		t.Fatalf("The content of the log file were wrong. expected: %s, %s, got: %s, %s", "INFO", testLogMessage, logMessages[0].Level, logMessages[0].Message)
	}
}
