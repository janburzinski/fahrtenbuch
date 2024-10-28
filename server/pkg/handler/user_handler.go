package handler

import (
	"net/http"
	"os"
	"server/pkg/logger"
	"server/pkg/middleware"
	"server/pkg/models"
	"server/pkg/repositories"
	"server/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTOkenDuration = 7 * 24 * time.Hour
)

var (
	AccessTokenSecret  = os.Getenv("ACCESS_TOKEN_SECRET")
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
)

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type RegisterRequestBody struct {
	FirstName string `form:"firstname"`
	LastName  string `form:"lastname"`
	Email     string `form:"email"`
	Password  string `form:"password"`
}

func Register(c *gin.Context) {
	var requestBody RegisterRequestBody
	if err := c.ShouldBind(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "wrong request body",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(requestBody.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error while hashing password",
		})
		logger.Log(logger.LOG_ERROR, "error while hashing password: %s", err.Error())
		return
	}

	user := models.User{
		FirstName: requestBody.FirstName,
		LastName:  requestBody.LastName,
		Email:     requestBody.Email,
		Password:  hashedPassword,
	}
	result := repositories.CreateUser(&user)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": result.Error.Error(),
		})
		logger.Log(logger.LOG_ERROR, "error while creating user: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "account successfully created",
	})
}

func generateToken(user models.User, duration time.Duration, secret string) (string, error) {
	claims := TokenClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(secret))
}

type LoginRequestBody struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func Login(c *gin.Context) {
	var requestBody LoginRequestBody
	if err := c.ShouldBind(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "wrong request body",
		})
		return
	}

	result := repositories.GetUserByEmail(requestBody.Email)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "wrong email or password", // most likely the email does not exist
		})
		logger.Log(logger.LOG_ERROR, "error while getting user: %s", result.Error.Error())
		return
	}

	if !utils.CheckPasswordHash(requestBody.Password, result.Result.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "wrong email or password",
		})
		return
	}

	//generate accesstoken
	accessToken, err := generateToken(*result.Result, AccessTokenDuration, AccessTokenSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error generating access token",
		})
		logger.Log(logger.LOG_ERROR, "error generating access token: %s", err.Error())
		return
	}

	//generate refresh token
	refreshToken, err := generateToken(*result.Result, RefreshTOkenDuration, RefreshTokenSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error generating refresh token",
		})
		logger.Log(logger.LOG_ERROR, "error generating access token: %s", err.Error())
		return
	}

	//store refresh token
	if err := repositories.StoreRefreshToken(result.Result.ID, refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error storing refresh token",
		})
		logger.Log(logger.LOG_ERROR, "error storing refresh token: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "successfully logged in",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

type RefreshTokenBodyRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken(c *gin.Context) {
	var requestBody RefreshTokenBodyRequest
	if err := c.ShouldBind(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	token, err := jwt.ParseWithClaims(requestBody.RefreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(RefreshTokenSecret), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid refresh token",
		})
		return
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid refresh token claims",
		})
		return
	}

	//verify refresh token in database
	valid := repositories.ValidateRefreshToken(claims.UserID, requestBody.RefreshToken)
	if !valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "refresh token has been revoked",
		})
		return
	}

	// get user from database
	result := repositories.GetUserByID(claims.UserID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error retrieving user",
		})
		return
	}

	newAccessToken, err := generateToken(*result.Result, AccessTokenDuration, AccessTokenSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "error generating new access token",
		})
		return
	}

	c.SetCookie("access_token", newAccessToken, int(AccessTokenDuration.Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "token refreshed successfully",
		"accessToken": newAccessToken,
	})
}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		userID, _, _ := middleware.GetUserFromContext(c)
		repositories.DeleteRefreshToken(userID, refreshToken)
	}

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "succesfully logged out",
	})
}
