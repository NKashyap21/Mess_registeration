package staff

import (
	"github.com/LambdaIITH/mess_registration/config"
	"gorm.io/gorm"
)

type ScanningController struct {
	DB *gorm.DB
}

func InitStaffController() *ScanningController {
	return &ScanningController{
		DB: config.GetDB(),
	}
}
