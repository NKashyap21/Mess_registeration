package middlewares

import (
	"os"
	"strings"

	"github.com/LambdaIITH/mess_registration/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("id", claims.ID)
		c.Next()
	}
}

// Alternative method to check for jsontoken in body
func AuthMiddlewareWithBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenRequest struct {
			JsonToken string `json:"jsontoken"`
		}

		// Try to get token from header first
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if authHeader != "" {
			tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
		} else {
			// Try to get token from body
			if err := c.ShouldBindJSON(&tokenRequest); err != nil || tokenRequest.JsonToken == "" {
				utils.UnauthorizedResponse(c, "Authorization token is required")
				c.Abort()
				return
			}
			tokenString = tokenRequest.JsonToken
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("id", claims.ID)
		c.Next()
	}
}

// Middleware for mess API key authentication
func MessAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		expectedKey := os.Getenv("MESS_API_KEY")

		if apiKey == "" || apiKey != expectedKey {
			utils.UnauthorizedResponse(c, "Invalid or missing API key")
			c.Abort()
			return
		}

		c.Next()
	}
}
