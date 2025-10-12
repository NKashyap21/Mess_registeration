package logs

import (
	"net/http"
	"strconv"
	"time"

	"github.com/LambdaIITH/mess_registration/services"
	"github.com/gin-gonic/gin"
)

type LogsController struct {
	loggerService *services.LoggerService
}

// GetLogsHandler handles GET /api/admin/logs
func (lc *LogsController) GetLogsHandler(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")
	userIDStr := c.Query("user_id")
	action := c.Query("action")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 1000 {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var userID *uint
	if userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			userIDVal := uint(uid)
			userID = &userIDVal
		}
	}

	var startDate, endDate *time.Time
	if startDateStr != "" {
		if sd, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &sd
		}
	}
	if endDateStr != "" {
		if ed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set end date to end of day
			endOfDay := ed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &endOfDay
		}
	}

	logs, total, err := lc.loggerService.GetLogs(limit, offset, userID, action, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":   logs,
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// GetUserActivityHandler handles GET /api/admin/logs/user/:user_id
func (lc *LogsController) GetUserActivityHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	limitStr := c.DefaultQuery("limit", "100")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 500 {
		limit = 100
	}

	logs, err := lc.loggerService.GetUserActivity(uint(userID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve user activity",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id": userID,
			"logs":    logs,
			"count":   len(logs),
		},
	})
}

// GetSystemLogsHandler handles GET /api/admin/logs/system
func (lc *LogsController) GetSystemLogsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 500 {
		limit = 100
	}

	logs, err := lc.loggerService.GetSystemLogs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve system logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":  logs,
			"count": len(logs),
		},
	})
}

// GetLogStatsHandler handles GET /api/admin/logs/stats
func (lc *LogsController) GetLogStatsHandler(c *gin.Context) {
	stats, err := lc.loggerService.GetLogStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve log statistics",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ExportLogsHandler handles GET /api/admin/logs/export
func (lc *LogsController) ExportLogsHandler(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid start_date format. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid end_date format. Use YYYY-MM-DD",
		})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	data, err := lc.loggerService.ExportLogs(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to export logs",
			"error":   err.Error(),
		})
		return
	}

	filename := "logs_" + startDateStr + "_to_" + endDateStr + ".json"
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", data)
}

// GetLogsByDateRangeHandler handles GET /api/admin/logs/range
func (lc *LogsController) GetLogsByDateRangeHandler(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid start_date format. Use YYYY-MM-DD",
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid end_date format. Use YYYY-MM-DD",
		})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	logs, err := lc.loggerService.GetLogsByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":       logs,
			"count":      len(logs),
			"start_date": startDateStr,
			"end_date":   endDateStr,
		},
	})
}
