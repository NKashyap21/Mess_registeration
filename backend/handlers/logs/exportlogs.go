package logs

import (
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// ExportLogsHandler handles GET /api/admin/logs/export
func (lc *LogsController) ExportLogsHandler(c *gin.Context) {
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
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	data, err := lc.loggerService.ExportLogs(startDate, endDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to export logs: "+err.Error())
		return
	}

	filename := "logs_" + startDateStr + "_to_" + endDateStr + ".json"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", data)
}
