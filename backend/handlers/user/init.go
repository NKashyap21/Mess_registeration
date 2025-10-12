package user

import (
	"github.com/LambdaIITH/mess_registration/config"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func InitUserController() *UserController {
	return &UserController{
		DB: config.GetDB(),
	}
}
