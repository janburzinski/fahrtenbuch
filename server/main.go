package main

import (
	"fahrtenbuch/pkg/db"
	"fahrtenbuch/pkg/redis"
	"fahrtenbuch/pkg/routes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

const (
	logFileName = "app.log"
	maxLogSize  = 100 * 1024 * 1024 // 100 mb
)

var (
	logFile     *os.File
	multiWriter io.Writer
	mu          sync.Mutex
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env variables")
	}
}

func rotateLogFile() error {
	log.Printf("Starting log rotation....")
	mu.Lock()
	defer mu.Unlock()

	fileInfo, err := os.Stat(logFileName)
	if err != nil {
		return err
	}

	if fileInfo.Size() > maxLogSize {
		logFile.Close()

		err = os.Rename(logFileName, fmt.Sprintf("%s.%s", logFileName, time.Now().Format("2006-01-02")))
		if err != nil {
			return err
		}

		newLogFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		logFile = newLogFile
		multiWriter = io.MultiWriter(os.Stdout, logFile)
	}
	log.Printf("Finished log rotation!")

	return nil
}

func scheduleLogRotation() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		duration := next.Sub(now)

		time.Sleep(duration)

		err := rotateLogFile()
		if err != nil {
			log.Fatalf("error rotation log file: %s", err)
		}
	}
}

func main() {
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		// ETag:                  true,
		DisableStartupMessage: true,
	})
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept",
	}))

	// init logger and write to file
	// Create a log file
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	multiWriter = io.MultiWriter(os.Stdout, logFile)

	// Configure Fiber logger to use the custom writer
	app.Use(logger.New(logger.Config{
		Format:        "[${time}] ${status} - ${latency} ${method} ${path}\n",
		Output:        multiWriter,
		DisableColors: false,
	}))

	// schedule log rotation
	go scheduleLogRotation()

	// init postgres connection
	db.Init()

	// init redis connection
	client := redis.Connect()
	rdb := redis.Initialize(client)
	if rdb == nil {
		log.Fatal("Error while connecting to redis!!!!")
		return
	}

	defer rdb.Close()

	// init api routes
	routes.SetupRoutes(app)

	// start web server
	serverPort := os.Getenv("SERVER_PORT")
	log.Printf("Trying to start the webserver in port: %s", serverPort)
	log.Fatal(app.Listen(":" + serverPort))
}
