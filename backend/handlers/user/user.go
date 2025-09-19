package user

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
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

func (u *UserController) GetUserInfoHandler(c *gin.Context) {
	userID := utils.ValidateSession(c)

	var user models.User
	err := u.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user info")
		return
	}
	var mess_name string

	switch user.Mess {
	case 1:
		mess_name = "Mess A LDH"
	case 2:
		mess_name = "Mess A UDH"
	case 3:
		mess_name = "Mess B LDH"
	case 4:
		mess_name = "Mess B UDH"
	default:
		mess_name = "No mess assigned"
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "User info fetched successfully",
		Data: map[string]any{
			"name":        user.Name,
			"roll_number": user.RollNo,
			"mess_id":     mess_name,
			"user_type":   user.Type,
		},
	})
}
