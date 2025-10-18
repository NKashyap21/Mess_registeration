package hosteloffice

import (
	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/services"
	"gorm.io/gorm"
)

type OfficeController struct {
	DB           *gorm.DB
	redisService *services.RedisMessService
}

func InitOfficeController() *OfficeController {
	return &OfficeController{
		DB:           config.GetDB(),
		redisService: services.NewRedisMessService(),
	}
}
