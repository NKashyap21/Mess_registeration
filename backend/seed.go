package main

import (
	"log"
	// "time"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
)

// SeedDatabase inserts default records if they don't exist
func SeedDatabase() {
	db := config.DB

	// Check and insert mess_registration_details
	var messCount int64
	if err := db.Model(&models.MessRegistrationDetails{}).Count(&messCount).Error; err != nil {
		log.Fatal("Failed to count mess_registration_details:", err)
	}
	if messCount == 0 {
		record := models.MessRegistrationDetails{
			// VegRegistrationStart: time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC),
			MessBLDHCapacity: 100,
			MessBUDHCapacity: 120,
		}
		if err := db.Create(&record).Error; err != nil {
			log.Fatal("Failed to insert default mess_registration_details:", err)
		}
		log.Println("Inserted default mess_registration_details")
	}

	// Check and insert a default user
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Fatal("Failed to count users:", err)
	}
	if userCount == 0 {
		user := models.User{
			Name:        "Test User MUQ",
			Email:       "es23btech11028@iith.ac.in",
			Phone:       nil,
			RollNo:      "es23btech11028",
			Mess:        0,
			Type:        0,
			CanRegister: true,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatal("Failed to insert default user:", err)
		}
		log.Println("Inserted default user")
	}
}
