package swap

import (
	"github.com/LambdaIITH/mess_registration/config"
	"gorm.io/gorm"
)

type SwapController struct {
	DB *gorm.DB
}

func InitSwapController() *SwapController {
	return &SwapController{
		DB: config.GetDB(),
	}
}
