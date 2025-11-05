package status

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func InitHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) CheckHealth(c *gin.Context) {
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Server is running",
		Data: map[string]interface{}{
			"status":  "healthy",
			"version": "1.0.0",
		},
	})
}
