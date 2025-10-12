package swap

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (sc *SwapController) GetAllSwapRequestsHandler(c *gin.Context) {
	userID := utils.ValidateSession(c)

	var swapRequests []models.SwapRequest
	err := sc.DB.Table("swap_requests").
		Select("swap_requests.*, users.email, users.name").
		Joins("join users on users.id = swap_requests.user_id").
		Where("swap_requests.user_id != ?", userID).
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

