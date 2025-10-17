package migrations

import (
	"log"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/db"
)

func MigrateDB() error {
	database := config.GetDB()

	log.Printf("Starting database migration...")
	err := database.AutoMigrate(
		&db.User{},
		&db.LoggerDetails{}, // Add logs table
		&db.MessRegistrationDetails{},
		&db.SwapRequest{},
	)

	if err != nil {
		log.Printf("Error occurred during database migration: %v", err)
		return err
	}

	log.Printf("Database migration completed successfully.")
	return nil
}
