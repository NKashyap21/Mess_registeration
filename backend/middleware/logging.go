package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Logger middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Color codes for different status ranges
		var statusColor, methodColor, resetColor string
		if gin.IsDebugging() {
			statusColor = getStatusColor(param.StatusCode)
			methodColor = getMethodColor(param.Method)
			resetColor = "\033[0m"
		}

		return fmt.Sprintf("%s[%s]%s | %s%-7s%s | %-50s | %s%3d%s | %13v | %s\n",
			"\033[36m", // Cyan for timestamp
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			resetColor,
			methodColor,
			param.Method,
			resetColor,
			param.Path,
			statusColor,
			param.StatusCode,
			resetColor,
			param.Latency,
			param.ErrorMessage,
		)
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
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("panic recovered: %s", err)
		}
		c.JSON(500, gin.H{
			"success": false,
			"message": "Internal server error",
			"error":   "Something went wrong",
		})
	})
}
