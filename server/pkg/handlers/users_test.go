package handlers

import (
	"fahrtenbuch/pkg/db"
	"fahrtenbuch/pkg/util"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func setupTest(t *testing.T) *fiber.App {
	//init db
	db.Init(true)

	// init web server
	app := fiber.New(fiber.Config{
		CaseSensitive:         true,
		DisableStartupMessage: true,
	})
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE,OPTIONS",
		ExposeHeaders: "Content-Type,Authorization,Accept",
	}))

	//setup routes
	userHandler := NewUserHandler()
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	return app
}

func TestRegisterRequest(t *testing.T) {
	var tests = []struct {
		description   string
		expectedError bool
		expected      string
		validEmail    bool
	}{
		{
			description:   "try registering with an valid email",
			expectedError: false,
			expected:      "user was successfully created",
			validEmail:    true,
		},
		{
			description:   "try registering with an already existant email",
			expectedError: true,
			expected:      "email is already being used",
			validEmail:    false,
		},
		{
			description:   "try registering with an invalid email",
			expectedError: true,
			expected:      "invalid email given",
		},
	}

	app := setupTest()
	route := "/register"

	validEmail := generateEmail(true)
	invalidEmail := generateEmail(false)

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {

			req := httptest.NewRequest("POST", route, nil)
		})
	}
}

func TestLoginAndMeRequest(t *testing.T) {

}

func generateEmail(validEmail bool) (email string) {
	miauError := "error generating email miau"
	if validEmail {
		randomString, err := util.GenerateRandomString(12)
		if err != nil {
			return miauError
		}
		email = fmt.Sprintf("%s@gmail.com", randomString)
	} else {
		randomString, err := util.GenerateRandomString(12)
		if err != nil {
			return miauError
		}
		email = randomString
	}
	return email
}
