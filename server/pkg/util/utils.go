package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"regexp"
)

var (
	// variable in the .env file that is true when the server runs on prod env
	IsProd = os.Getenv("IS_PROD") == "true"
	// used for refresh token cookies
	MaxCookieAge = 1000 * 60 * 60 * 24 * 365 * 10 // 10 years
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomBytes(length uint32) ([]byte, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}
	return string(result), nil
}

func GenerateRandomInt(min, max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max-min+1))
	if err != nil {
		panic(err)
	}
	return nBig.Int64() + min
}

func GenerateRandomPhoneNumber() string {
	areaCode := GenerateRandomInt(1, 999)
	rest := GenerateRandomInt(1, 999)

	return fmt.Sprintf("(%03d) %03d-%04d", areaCode, rest, rest)
}

func ValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
