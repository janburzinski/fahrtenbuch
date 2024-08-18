package util

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

// time is supposed to be a given in unix
func CreateToken(userId string, accessToken bool) (string, error) {
	accessTokenJwtTokenExp := time.Now().Add(time.Minute * 15).Unix()

	refreshTokenJwtTokenExp := time.Now().Add(time.Hour * 720).Unix() // 1 month

	switch accessToken {
	case true:
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"userId": userId,
				"exp":    accessTokenJwtTokenExp,
			})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return "", err
		}

		return tokenString, nil

	default:
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"userId": userId,
				"exp":    refreshTokenJwtTokenExp,
			})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid access jwt token")
	}

	return nil
}

func GetUserIdFromJWT(tokenString string) (userId string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := fmt.Sprint(claims["userId"])
		if userId == "" {
			return "", errors.New("invalid jwt token")
		}
		return userId, nil
	} else {
		return "", err
	}
}
