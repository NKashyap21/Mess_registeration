package hosteloffice

import (
	"log"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (oc *OfficeController) ApplyNewRegistration(c *gin.Context) {
	// Step 1: Ensure at least one row exists
	var reg models.MessRegistrationDetails
	if err := oc.DB.First(&reg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			reg = models.MessRegistrationDetails{
				VegRegistrationOpen:    false,
				NormalRegistrationOpen: false,
			}
			if err := oc.DB.Create(&reg).Error; err != nil {
				utils.RespondWithError(c, http.StatusInternalServerError, "Failed to initialize registration details")
				return
			}
		} else {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch registration details")
			return
		}
	}

	// Step 2: Close both registrations globally
	if err := oc.DB.Model(&models.MessRegistrationDetails{}).
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Updates(map[string]interface{}{
			"veg_registration_open":    false,
			"normal_registration_open": false,
		}).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to close registrations")
		return
	}

	// Step 3: Start a transaction
	tx := oc.DB.Begin()
	if tx.Error != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	// Step 4: Copy next_mess → mess
	if err := tx.Model(&models.User{}).
		Where("next_mess IS NOT NULL AND can_register = ?", true).
		Updates(map[string]interface{}{
			"mess": gorm.Expr("next_mess"),
		}).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to copy next_mess to mess")
		return
	}

	// Step 5: Reset next_mess to 0 (or NULL if preferred)
	if err := tx.Model(&models.User{}).Updates(map[string]interface{}{"next_mess": 0}).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to reset next_mess")
		return
	}

	// Step 6: Commit
	if err := tx.Commit().Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	if err := oc.redisService.ClearMessCount(); err != nil {
		log.Println("Failed to clear redis mess counts")
		log.Println(err)
	}

	// Step 7: Return updated registration details
	if err := oc.DB.First(&reg).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch updated registration status")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "New registration cycle applied successfully",
		Data: map[string]interface{}{
			"veg_registration_open":    reg.VegRegistrationOpen,
			"normal_registration_open": reg.NormalRegistrationOpen,
		},
	})
}
