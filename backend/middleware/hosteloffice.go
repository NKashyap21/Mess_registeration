package middleware

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HostelOfficeMiddleWare(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check for user_id
		user_id, ok := c.Request.Context().Value("user_id").(uint)
		if !ok || user_id == 0 {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid User Context.")
			c.Abort()
			return
		}
		//Check if the person Exists in the Database
		var user models.User
		if err := db.First(&user, user_id).Error; err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}
		//Check if the person is Hostel Office.
		if user.Type != 2 {
			utils.RespondWithError(c, http.StatusForbidden, "Access Denied: Hostel Office Only.")
			c.Abort()
			return
		}

		c.Next()
	}
}

func MessStaffMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check for user_id
		user_id, ok := c.Request.Context().Value("user_id").(uint)
		if !ok || user_id == 0 {
			utils.RespondWithError(c, http.StatusUnauthorized, "Invalid User Context.")
			c.Abort()
			return
		}
		//Check if the person Exists in the Database
		var user models.User
		if err := db.First(&user, user_id).Error; err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}
		//Check if the person is Mess Staff.
		if user.Type != 1 {
			utils.RespondWithError(c, http.StatusForbidden, "Access Denied: Mess Staff Only.")
			c.Abort()
			return
		}

		// Add user to context for use in handlers
		c.Set("user", user)
		c.Next()
	}
}
