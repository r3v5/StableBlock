package utils

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
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}).SignedString(GetSecret())

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(time.Minute * 10).Unix(),
	}).SignedString(GetSecret())

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}


func GenerateAccessToken(address string) (string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}).SignedString(GetSecret())

	return accessToken, err
}
