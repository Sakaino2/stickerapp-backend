package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, userId uint) (string, string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	}).SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 10).Unix(),
	}).SignedString([]byte(os.Getenv("SECRET")))

	return token, refreshToken, err
}

func VerifyToken(token string) (int64, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return 0, errors.New("could not parse token")
	}
	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Claims not ok")
		return 0, errors.New("invalid token claims")
	}

	//email, _ := claims["email"].(string)
	userId, _ := claims["userId"].(float64)
	return int64(userId), nil
}

func RefreshToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", "", err
	}

	email, ok := claims["email"].(string)

	if !ok {
		return "", "", err
	}

	userId, ok := claims["userId"].(string)

	if !ok {
		return "", "", err
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	})

	newTokenString, err := newToken.SignedString(os.Getenv("SECRET"))

	if err != nil {
		return "", "", err
	}

	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 15).Unix(), // Refresh token expires in 7 days
	})

	newRefreshTokenString, err := newRefreshToken.SignedString(os.Getenv("SECRET")) // Replace "your-secret-key" with your own secret key

	if err != nil {
		return "", "", err
	}

	return newTokenString, newRefreshTokenString, nil
}
