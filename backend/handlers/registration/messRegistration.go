package registration

import (
	"net/http"
	"strconv"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MessController struct {
	DB *gorm.DB
}

func InitMessController() *MessController {
	return &MessController{DB: config.GetDB()}
}

func (m *MessController) MessRegistrationHandler(c *gin.Context) {
	// Only accept requests on this endpoint at a specified date
	// Check if the current date is within the registration period
	if !isRegistrationOpen() {
		utils.RespondWithError(c, http.StatusForbidden, "Registration is not open at this time")
		return
	}

	userID := utils.ValidateSession(c)

	// Check if user already has a mess assigned
	var user models.User
	if err := m.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		return
	}

	if !user.CanRegister {
		utils.RespondWithError(c, http.StatusBadRequest, "User cannot Register.")
		return
	}

	if user.Mess != 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "User already has a mess assigned")
		return
	}

	messParam := c.Param("mess")
	mess, err := strconv.Atoi(messParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Mess parameter must be an integer")
		return
	}

	if !isValidMess(mess) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid parameter, mess must be between 1 and 4")
		return
	}

	// Set user's mess to the specified mess
	user.Mess = int8(mess)
	if err := m.DB.Save(&user).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update mess: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Mess registration successful",
	})

}

func isValidMess(mess int) bool {
	return mess >= 1 && mess < 5
}

func isRegistrationOpen() bool {
	return true // Placeholder: Implement actual date check logic
}
