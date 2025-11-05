package hosteloffice

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (oc *OfficeController) ToggleNormalRegistration(c *gin.Context) {
	// Toggle NormalRegistrationOpen for the first row
	if err := oc.DB.Model(&db.MessRegistrationDetails{}).Session(&gorm.Session{AllowGlobalUpdate: true}).
		Update("normal_registration_open", gorm.Expr("NOT normal_registration_open")).
		Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to toggle normal registration")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Normal registration toggled successfully",
	})
}

func (oc *OfficeController) ToggleVegRegistration(c *gin.Context) {
	// Toggle VegRegistrationOpen for the first row
	if err := oc.DB.Model(&db.MessRegistrationDetails{}).Session(&gorm.Session{AllowGlobalUpdate: true}).
		Update("veg_registration_open", gorm.Expr("NOT veg_registration_open")).
		Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to toggle veg registration")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Veg registration toggled successfully",
	})
}
