package utils

import (
	"github.com/LambdaIITH/mess_registration/db"
	"gorm.io/gorm"
)

func GetNormalRegistrationStatus(DB *gorm.DB) bool {
	var registrationDetails db.MessRegistrationDetails
	if err := DB.First(&registrationDetails).Error; err != nil {
		return false
	}

	return registrationDetails.NormalRegistrationOpen
}


func GetVegRegistrationStatus(DB *gorm.DB) bool {
	var registrationDetails db.MessRegistrationDetails
	if err := DB.First(&registrationDetails).Error; err != nil {
		return false
	}

	return registrationDetails.VegRegistrationOpen
}
