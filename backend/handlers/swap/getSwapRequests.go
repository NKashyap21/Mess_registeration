package swap

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (sc *SwapController) GetAllSwapRequestsHandler(c *gin.Context) {
	// userID := utils.ValidateSession(c)

	var swapRequests []models.SwapRequest
	err := sc.DB.Table("swap_requests").
		Select("swap_requests.*, users.email, users.name").
		Joins("join users on users.id = swap_requests.user_id").
		// Where("swap_requests.user_id != ?", userID).
		Find(&swapRequests).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch swap requests")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Swap requests fetched successfully",
		Data:    swapRequests,
	})
}

func (sc *SwapController) GetSwapRequestsByID(c *gin.Context) {
	userID := utils.ValidateSession(c)

	var swapRequest db.SwapRequest
	err := sc.DB.Table("swap_requests").
		Select("swap_requests.*, users.email, users.name").
		Joins("join users on users.id = swap_requests.user_id").
		Where("swap_requests.user_id = ?", userID).
		First(&swapRequest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "No swap request found for the user")
		} else {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch swap request")
		}
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Swap request fetched successfully",
		Data:    swapRequest,
	})
}
