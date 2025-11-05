package logs

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// GetLogStatsHandler handles GET /api/admin/logs/stats
func (lc *LogsController) GetLogStatsHandler(c *gin.Context) {
	stats, err := lc.loggerService.GetLogStats()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve log statistics: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: stats,
	})
}
