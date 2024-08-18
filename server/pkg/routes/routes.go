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
	tokenString := c.Cookies("jwt_token")

	// check accessToken
	err := util.VerifyToken(tokenString)
	if err != nil {
		// access token might just be expired
		// check refresh Token
		// and possibly refresh access token
		refreshToken := c.Cookies("jwt_refresh_token")
		err := util.VerifyToken(refreshToken)
		if err != nil {
			errorResp := handlers.ErrorResponse{
				OK:    false,
				Error: "not logged in",
			}
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}

		// access token just expired
		// issue a new one
		userId, err := util.GetUserIdFromJWT(refreshToken)
		if err != nil {
			errorResp := handlers.ErrorResponse{
				OK:    false,
				Error: "not logged in",
			}
			return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
		}
		accessToken, err := util.CreateToken(userId, true)
		if err != nil {
			errorResp := handlers.ErrorResponse{
				OK:    false,
				Error: "error while creating access token",
			}
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt_token",
			Value:    accessToken,
			Expires:  time.Now().Add(time.Minute * 15),
			HTTPOnly: util.IsProd,
			Secure:   true,
			SameSite: fiber.CookieSameSiteLaxMode,
			MaxAge:   15 * 60, // 15 min
		})

		// also refresh token (does this even make sense?)
		newRefreshToken, err := util.CreateToken(userId, false)
		if err != nil {
			errorResp := handlers.ErrorResponse{
				OK:    false,
				Error: "error while creating refresh token",
			}
			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt_refresh_token",
			Value:    newRefreshToken,
			Expires:  time.Now().Add(time.Hour * 720),
			HTTPOnly: util.IsProd,
			Secure:   true,
			SameSite: fiber.CookieSameSiteLaxMode,
			MaxAge:   util.MaxCookieAge, // 10 years
		})

		return c.Next()
	}

	// access token still valid
	// all good
	return c.Next()
}
