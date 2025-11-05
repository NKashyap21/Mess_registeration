package logs

import (
	"net/http"
	"strconv"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// GetUserActivityHandler handles GET /api/admin/logs/user/:user_id
func (lc *LogsController) GetUserActivityHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	limitStr := c.DefaultQuery("limit", "100")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 500 {
		limit = 100
	}

	logs, err := lc.loggerService.GetUserActivity(uint(userID), limit)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve user activity: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"user_id": userID,
			"logs":    logs,
			"count":   len(logs),
		},
	})
}
