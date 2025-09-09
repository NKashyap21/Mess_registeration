package middlewares

import (
	"github.com/LambdaIITH/mess_registration/internal/services"
	"github.com/LambdaIITH/mess_registration/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AdminMiddleware checks if the authenticated user is an admin
func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		id, exists := c.Get("id")
		if !exists {
			utils.UnauthorizedResponse(c, "User not authenticated")
			c.Abort()
			return
		}

		userID := id.(uuid.UUID)
		userService := services.NewUserService(db)

		user, err := userService.GetUserByID(userID)
		if err != nil {
			utils.UnauthorizedResponse(c, "User not found")
			c.Abort()
			return
		}

		// Check if user is admin
		if user.UserType != "admin" {
			utils.ErrorResponse(c, 403, "Forbidden", "Admin access required")
			c.Abort()
			return
		}

		// Set user info in context for use in handlers
		c.Set("user", user)
		c.Next()
	}
}

// MessStaffMiddleware checks if the authenticated user is mess staff
func MessStaffMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		id, exists := c.Get("id")
		if !exists {
			utils.UnauthorizedResponse(c, "User not authenticated")
			c.Abort()
			return
		}

		userID := id.(uuid.UUID)
		userService := services.NewUserService(db)

		user, err := userService.GetUserByID(userID)
		if err != nil {
			utils.UnauthorizedResponse(c, "User not found")
			c.Abort()
			return
		}

		// Check if user is mess staff or admin
		if user.UserType != "mess_staff" && user.UserType != "admin" {
			utils.ErrorResponse(c, 403, "Forbidden", "Mess staff access required")
			c.Abort()
			return
		}

		// Set user info in context for use in handlers
		c.Set("user", user)
		c.Next()
	}
}
