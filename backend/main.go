package main

import (
	"log"
	"os"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/internal/db"
	"github.com/LambdaIITH/mess_registration/internal/router"
	"github.com/gin-gonic/gin"
)

// func init() {
// 	config.LoadEnvVaariables()
// 	// Load config.json
// 	config.LoadConfig()

// 	// Connect DB
// 	config.ConnectDB()
// }

func main() {
	// Load configuration
	config.LoadEnvVariables()

	//Initialize database
	database, err := db.InitDB()

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.MigrateDB(database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	gin.SetMode(os.Getenv("GIN_MODE"))

	r := router.SetupRouter(database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on port %s", port)
	if r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
