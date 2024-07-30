package routes

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

const (
	rateLimitMax = 30
	currVersion  = "v1"
)

/*
*	Init all routes
 */
func SetupRoutes(app *fiber.App) {
	apiAuth := app.Group("/api/" + currVersion + "/auth")
	setupAuthRoutes(apiAuth)

	// setup rate limiter
	isProd := os.Getenv("IS_PROD") == "true"
	if isProd {
		app.Use(limiter.New(limiter.Config{
			Expiration: 10 * time.Second,
			Max:        rateLimitMax,
		}))
	}

	log.Print("Initialized Routes")
}

func setupAuthRoutes(app fiber.Router) {
}
