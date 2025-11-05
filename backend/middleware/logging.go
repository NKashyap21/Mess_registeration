package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/services"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// Logger middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	logger := services.GetLoggerService()

	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered interface{}) {
		// Handle panic recovery
		if err, ok := recovered.(string); ok {
			log.Printf("panic recovered: %s", err)
			logger.LogSystemAction("PANIC_RECOVERY", fmt.Sprintf("Panic recovered: %s", err))
		}
		utils.RespondWithJSON(c, http.StatusInternalServerError, models.APIResponse{
			Message: "Internal server error",
		})
	})
}

// DatabaseLogger middleware specifically for logging HTTP requests to database
func DatabaseLogger() gin.HandlerFunc {
	logger := services.GetLoggerService()

	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		end := time.Now()
		latency := end.Sub(start)

		// Get user ID from context if available
		var userID uint = 0
		if userInterface, exists := c.Get("user_id"); exists {
			if uid, ok := userInterface.(uint); ok {
				userID = uid
			} else if uidStr, ok := userInterface.(string); ok {
				if uid, err := strconv.ParseUint(uidStr, 10, 32); err == nil {
					userID = uint(uid)
				}
			}
		}

		// Prepare full path
		if raw != "" {
			path = path + "?" + raw
		}

		// Create log message
		message := fmt.Sprintf("%s %s - %d - %v",
			c.Request.Method,
			path,
			c.Writer.Status(),
			latency)

		// Log to database
		logger.LogHTTPRequest(c, userID, c.Writer.Status(), message)

		// Also log to console with colors for development
		if gin.IsDebugging() {
			statusColor := getStatusColor(c.Writer.Status())
			methodColor := getMethodColor(c.Request.Method)
			resetColor := "\033[0m"

			fmt.Printf("%s[%s]%s | %s%-7s%s | %-50s | %s%3d%s | %13v | %s\n",
				"\033[36m", // Cyan for timestamp
				end.Format("2006/01/02 - 15:04:05"),
				resetColor,
				methodColor,
				c.Request.Method,
				resetColor,
				path,
				statusColor,
				c.Writer.Status(),
				resetColor,
				latency,
				c.Errors.ByType(gin.ErrorTypePublic).String(),
			)
		}
	})
}

func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" // Green
	case code >= 300 && code < 400:
		return "\033[33m" // Yellow
	case code >= 400 && code < 500:
		return "\033[31m" // Red
	case code >= 500:
		return "\033[35m" // Magenta
	default:
		return "\033[37m" // White
	}
}

func getMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Blue
	case "POST":
		return "\033[32m" // Green
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[35m" // Magenta
	case "HEAD":
		return "\033[36m" // Cyan
	case "OPTIONS":
		return "\033[37m" // White
	default:
		return "\033[0m" // Reset
	}
}

// Recovery middleware recovers from any panics and writes a 500 if there was one
func Recovery() gin.HandlerFunc {
	logger := services.GetLoggerService()

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("panic recovered: %s", err)
			logger.LogSystemAction("PANIC_RECOVERY", fmt.Sprintf("Panic recovered at %s %s: %s", c.Request.Method, c.Request.URL.Path, err))
		} else {
			log.Printf("panic recovered: %v", recovered)
			logger.LogSystemAction("PANIC_RECOVERY", fmt.Sprintf("Panic recovered at %s %s: %v", c.Request.Method, c.Request.URL.Path, recovered))
		}
		utils.RespondWithJSON(c, http.StatusInternalServerError, models.APIResponse{
			Message: "Internal server error",
		})
	})
}
