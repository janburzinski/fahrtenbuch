package routes

import (
	"fahrtenbuch/pkg/handlers"
	"fahrtenbuch/pkg/util"
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

	apiUser := app.Group("/api/"+currVersion+"/user", authMiddleware)
	setupUserRoutes(apiUser)

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

func authMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("jwt_secret")

	// check accessToken
	err := util.VerifyToken(tokenString, true)
	if err != nil {
		errorResp := handlers.ErrorResponse{
			OK:    false,
			Error: "not logged in",
		}
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	// check refresh Token (refresh?)

	// verify with user id from redis

	return c.Next()
}
