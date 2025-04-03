package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func GetSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateTokens(address string) (string, string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(time.Minute * 1).Unix(),
	}).SignedString(GetSecret())

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(time.Minute * 3).Unix(),
	}).SignedString(GetSecret())

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
