package main

import (
	"log"
	"os"
	"strconv"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/migrations"
	"github.com/LambdaIITH/mess_registration/router"
	"github.com/LambdaIITH/mess_registration/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	config.ConnectDatabase()

	// Initialize Redis connection
	config.ConnectRedis()

	// Run database migrations
	if err := migrations.MigrateDB(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize Logger Service
	loggerService := services.GetLoggerService()
	loggerService.LogSystemAction("SYSTEM_START", "Application starting up")

	// Setup graceful shutdown for logger service
	defer func() {
		loggerService.LogSystemAction("SYSTEM_SHUTDOWN", "Application shutting down")
		loggerService.Shutdown()
	}()

	// Initialize sync service and sync Redis from DB
	syncService := services.NewSyncService()
	if err := syncService.InitializeRedisFromDB(); err != nil {
		log.Fatal("Failed to initialize Redis from database:", err)
	}

	// Start background sync service (will only sync when registration is open)
	syncIntervalStr := os.Getenv("SYNC_INTERVAL_SECONDS")
	syncInterval := 30 // default 30 seconds
	if syncIntervalStr != "" {
		if interval, err := strconv.Atoi(syncIntervalStr); err == nil {
			syncInterval = interval
		}
	}
	log.Println("Starting sync service - will only sync when registration is open")
	syncService.StartBackgroundSync(syncInterval)

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
