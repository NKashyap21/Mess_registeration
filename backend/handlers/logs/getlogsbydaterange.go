package logs

import (
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// GetLogsByDateRangeHandler handles GET /api/admin/logs/range
func (lc *LogsController) GetLogsByDateRangeHandler(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid to_date format. Use YYYY-MM-DD")
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	logs, err := lc.loggerService.GetLogsByDateRange(startDate, endDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve logs: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"logs":       logs,
			"count":      len(logs),
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	})
}
