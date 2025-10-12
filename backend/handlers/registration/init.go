package registration

import (
	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/services"
	"gorm.io/gorm"
)

type MessController struct {
	DB           *gorm.DB
	redisService *services.RedisMessService
}

func InitMessController() *MessController {
	return &MessController{
		DB:           config.GetDB(),
		redisService: services.NewRedisMessService(),
	}
}
