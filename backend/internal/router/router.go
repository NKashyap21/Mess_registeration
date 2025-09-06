package router

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Custom log formatter (like in gian_portal)
func LogFormatter(params gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] - %s \"%s %s %s %d %s \"%s\" %s\"\n",
		params.TimeStamp.Format("2006-01-02 15:04:05"),
		params.ClientIP,
		params.Method,
		params.Path,
		params.Request.Proto,
		params.StatusCode,
		params.Latency,
		params.Request.UserAgent(),
		params.ErrorMessage,
	)
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	fmt.Println("\033[36mMess Registration server started.\033[0m")

	// Custom logger
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output:    os.Stdout,
		Formatter: LogFormatter,
	}))

	// Optional: add request logging middleware if you create one
	// router.Use(middlewares.RequestLoggerMiddleware)

	// ✅ CORS setup
	config := cors.Config{
		AllowOrigins:     []string{os.Getenv("WEB_URL")}, // must include http:// or https://
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))

	// Routes
	SetupRoutes(router)

	return router
}
