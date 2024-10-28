package middleware

import (
	"errors"
	"net/http"
	"server/pkg/handler"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrNoAuthHeader      = errors.New("no authorization header")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		claims, err := validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid or expired token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func extractToken(c *gin.Context) (token string, err error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return "", ErrInvalidAuthFormat
		}
		return parts[1], nil
	}

	token, err = c.Cookie("access_token")
	if err != nil {
		return "", ErrNoAuthHeader
	}
	return token, nil
}

func validateToken(tokenString string) (*handler.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &handler.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(handler.AccessTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*handler.TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func GetUserFromContext(c *gin.Context) (uint, string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, "", errors.New("user id not found in context")
	}

	email, exists := c.Get("email")
	if !exists {
		return 0, "", errors.New("email not found in context")
	}

	return userID.(uint), email.(string), nil
}
