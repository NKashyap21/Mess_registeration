package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, email string, name string, picture string, secretKey string) (string, error) {

	expirationTime := time.Now().Add(72 * time.Hour)

	claims := &jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"picture": picture,
		"name":    name,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "a", err
	}

	return tokenString, nil

}
