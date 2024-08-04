package main

import (
	"fahrtenbuch/pkg/db"
	"fahrtenbuch/pkg/redis"
	"fahrtenbuch/pkg/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env variables")
	}
}

func main() {
	app := fiber.New()
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())

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
	log.Fatal(app.Listen(":" + serverPort))
}
