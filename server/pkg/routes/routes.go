package routes

import (
	"fahrtenbuch/pkg/handlers"
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
**	Init all routes
 */
func SetupRoutes(app *fiber.App) {
	apiAuth := app.Group("/api/" + currVersion + "/auth")
	setupAuthRoutes(apiAuth)

	apiUser := app.Group("/api/"+currVersion+"/user", authMiddleware)
	setupUserRoutes(apiUser)

	apiOrganisation := app.Group("/api/"+currVersion+"/org", authMiddleware)
	setupOrganisationRoutes(apiOrganisation)

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
	userHandler := handlers.NewUserHandler()
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)
}

func setupUserRoutes(app fiber.Router) {
	// routes
	userHandler := handlers.NewUserHandler()
	app.Get("/me", userHandler.Me)
}

func setupOrganisationRoutes(app fiber.Router) {

}

func authMiddleware(c *fiber.Ctx) error {
	//implement firebase auth middleware
	return c.Next()
}
