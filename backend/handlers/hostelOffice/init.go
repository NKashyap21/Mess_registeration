package hosteloffice

import (
	"github.com/LambdaIITH/mess_registration/config"
	"gorm.io/gorm"
)

type OfficeController struct {
	DB *gorm.DB
}

func InitOfficeController() *OfficeController {
	return &OfficeController{
		DB: config.GetDB(),
	}
}
