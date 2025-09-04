package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func home(c *gin.Context) {
	year := time.Now().Year()
	c.String(http.StatusOK, "Mess Registration ©%d IIT Hyderabad", year)
}

func SetupRoutes(router *gin.Engine) {
	// Health Check
	router.GET("/", home)

	// Example API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "mess_registration",
		})
	})

	// Later you’ll add mess-specific routes like:
	// router.POST("/register", controller.RegisterForMess)
	// router.GET("/menu", controller.GetMenu)
}
