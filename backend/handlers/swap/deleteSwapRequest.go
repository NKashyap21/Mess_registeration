package swap

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func (sc *SwapController) DeleteSwapHandler(c *gin.Context) {
	userID := utils.ValidateSession(c)
	
	// Check if there exists a swap request for this user
	var existingRequest db.SwapRequest
	err := sc.DB.First(&existingRequest, "user_id = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		utils.RespondWithError(c, http.StatusBadRequest, "No swap request found for this user")
		return
	} else if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	// Delete the existing swap request
	err = sc.DB.Delete(&existingRequest, "user_id = ?", userID).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete swap request")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Swap request deleted successfully",
	})
}

