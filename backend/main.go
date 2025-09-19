package main

import (
	"log"
	"os"

	"github.com/LambdaIITH/mess_registration/config"
	// "github.com/LambdaIITH/mess_registration/migrations"
	"github.com/LambdaIITH/mess_registration/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	config.ConnectDatabase()

	// Run database migrations
	// if err := migrations.MigrateDB(); err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }

	// Setup router
	r := router.SetupRouter()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
