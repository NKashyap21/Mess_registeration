package migrations

import (
	"log"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
)

func MigrateDB() error {
	db := config.GetDB()
	
	log.Printf("Starting database migration...")
	err := db.AutoMigrate(
		&models.User{},
	)
	
	if err != nil {
		log.Printf("Error occurred during database migration: %v", err)
		return err
	}

	log.Printf("Database migration completed successfully.")
	return nil
}
