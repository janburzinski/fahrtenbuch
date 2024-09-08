package handlers

import (
	"fahrtenbuch/pkg/db"
	"fahrtenbuch/pkg/models"
	"fahrtenbuch/pkg/util"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		errResp := ErrorResponse{
			OK:    false,
			Error: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// check if the given email is valid
	valid := util.ValidEmail(user.Email)
	if !valid {
		errResp := ErrorResponse{
			OK:    false,
			Error: "invalid email",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// check if email already exists
	// if yes, send custom error message
	var count int64
	if db.DB.Model(&models.User{}).Where("email = ?", user.Email).Count(&count); count > 0 {
		errResp := ErrorResponse{
			OK:    false,
			Error: "email is already being used",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errResp)
	}

	// all good, so hash the password
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
		errResp := ErrorResponse{
			OK:    false,
			Error: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}
	user.Password = hashedPassword

	// create the user
	if err := db.DB.Create(&user).Error; err != nil {
		errResp := ErrorResponse{
			OK:    false,
			Error: err.Error(),
		}
		return c.Status(fiber.StatusInternalServerError).JSON(errResp)
	}

	resp := SuccessResponse{
		OK:      true,
		Message: "user was successfully created",
	}
	return c.JSON(resp)
}

func (uh *UserHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		errorResp := ErrorResponse{
			OK:    false,
			Error: err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var user models.User
	result := db.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		errorResp := ErrorResponse{
			OK:    false,
			Error: "user not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	match, err := util.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil || !match {
		errorResp := ErrorResponse{
			OK:    false,
			Error: "incorrect login data",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	// generate jwt and set as cookie / local storage
	userId := fmt.Sprintf("%d", user.ID)
	accessToken, err := util.CreateToken(userId, true)
	if err != nil {
		errorResp := ErrorResponse{
			OK:    false,
			Error: "error while creating access token",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	// set access token
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: util.IsProd,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   15 * 60, // 15 minutes
	})

	// generate refresh token
	refreshToken, err := util.CreateToken(userId, false)
	if err != nil {
		errorResp := ErrorResponse{
			OK:    false,
			Error: "error while creating refresh token",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	// set jwt refresh token as a cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt_refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 720),
		HTTPOnly: util.IsProd,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
		MaxAge:   util.MaxCookieAge, // 10 years
	})

	successResp := SuccessResponse{
		OK:      true,
		Message: "successfully logged in",
	}
	return c.JSON(successResp)
}

func (uh *UserHandler) Me(c *fiber.Ctx) error {
	token := c.Cookies("jwt_refresh_token") // get user id from jwt token (refreshtoken => longer ttl)
	userId, err := util.GetUserIdFromJWT(token)
	if err != nil {
		errorResp := ErrorResponse{
			OK:    true,
			Error: "user is not logged in",
		}
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	var user models.User
	result := db.DB.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		errorResp := ErrorResponse{
			OK:    false,
			Error: "user not found",
		}
		return c.Status(fiber.StatusNotFound).JSON(errorResp)
	}

	successResp := SuccessUserMeResponse{
		OK:      true,
		Message: "success",
		User:    user,
	}
	return c.JSON(successResp)
}

func (uh *UserHandler) VerifyToken(c *fiber.Ctx) error {
	return nil
}
