package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
)

func TokenRequired(db *gorm.DB, c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Fallback: try to get token from Authorization header for Safari compatibility
		// authHeader := c.Request.Header.Get("Authorization")
		// if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		// 	tokenString = authHeader[7:]
		// }
		tokenString, err := c.Cookie("jwt")

		if err == http.ErrNoCookie {
			utils.RespondWithError(c, http.StatusUnauthorized, "Token is missing!")
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJWTConfig().SecretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token signature")
				return
			}
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		if !token.Valid {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID in token")
			return
		}
		userID := uint(userIDFloat)

		// Verify user exists
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "User not found")
			return
		}

		// Add user ID to the context for the next handler
		ctx := context.WithValue(c.Request.Context(), "user_id", userID) //lint:ignore SA1029 ignore context key naming
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
