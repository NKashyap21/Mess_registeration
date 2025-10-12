package auth

import (
	"github.com/LambdaIITH/mess_registration/config"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func InitAuthController() *AuthController {
	return &AuthController{
		DB: config.GetDB(),
	}
}
