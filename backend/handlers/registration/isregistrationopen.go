package registration

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (m *MessController) IsRegistrationOpen(c *gin.Context) {
	utils.RespondWithJSON(c, http.StatusOK, map[string]bool{
		"regular": m.isRegistrationOpen(),
		"veg":     m.isVegRegistrationOpen(),
		// "regular": state.GetRegistrationStatusReg(),
		// "veg":     state.GetRegistrationStatusVeg(),
	})
}

func (m *MessController) isRegistrationOpen() bool {
	// Check if normal registration is open from the database

	var registrationDetails db.MessRegistrationDetails
	if err := m.DB.First(&registrationDetails).Error; err != nil {
		return false
	}

	return registrationDetails.NormalRegistrationOpen
}

func (m *MessController) isVegRegistrationOpen() bool {
	// Check if veg registration is open from the database

	var registrationDetails db.MessRegistrationDetails
	if err := m.DB.First(&registrationDetails).Error; err != nil {
		return false
	}

	return registrationDetails.VegRegistrationOpen
}
