package status

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func InitHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Server is running",
		"data": gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		},
	})
}

