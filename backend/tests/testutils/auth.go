package testutils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(email string, userID uint, userType int8) string {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_type": userType,
		"email":     email,
		"name":      "Test User",
		"picture":   "",
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := t.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return "Bearer " + s
}
