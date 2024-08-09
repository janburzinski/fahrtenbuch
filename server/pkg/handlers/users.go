package handlers

import (
	"fahrtenbuch/pkg/db"
	"fahrtenbuch/pkg/models"
	"fahrtenbuch/pkg/util"

	"github.com/gofiber/fiber/v2"
)

// use for dependency injection etc. later on
type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error()) // update the way the error is being displayed (send json instead of normal string)
	}

	p := &util.Argon2Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	// hash the password
	hashedPassword, err := util.GenerateFromPassword(user.Password, p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	user.Password = hashedPassword

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(user)
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var user models.User
	result := db.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).SendString("user not found")
	}

	match, err := util.ComparePasswordAndHash(user.Password, input.Password)
	if err != nil || !match {
		return c.Status(fiber.StatusBadRequest).SendString("incorrect login data")
	}

	// generate jwt and set as cookie

	return c.SendString("successful login")
}
