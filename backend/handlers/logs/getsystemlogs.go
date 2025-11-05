package logs

import (
	"net/http"
	"strconv"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// GetSystemLogsHandler handles GET /api/admin/logs/system
func (lc *LogsController) GetSystemLogsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 500 {
		limit = 100
	}

	logs, err := lc.loggerService.GetSystemLogs(limit)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve system logs: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"logs":  logs,
			"count": len(logs),
		},
	})
}
