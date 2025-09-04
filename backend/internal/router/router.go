package router

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	fmt.Println("\033[36mMess Registration server started.\033[0m")

	// Logger
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS setup
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{os.Getenv("WEB_URL")}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"X-Requested-With", "Content-Type", "Accept"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Routes
	SetupRoutes(router)

	return router
}
