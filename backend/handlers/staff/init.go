package staff

import (
	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/services"
	"gorm.io/gorm"
)

type ScanningController struct {
	DB           *gorm.DB
	RedisService *services.RedisMessService
}

func InitStaffController() *ScanningController {
	return &ScanningController{
		DB:           config.GetDB(),
		RedisService: services.NewRedisMessService(),
	}
}
