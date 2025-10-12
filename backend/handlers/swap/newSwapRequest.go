package swap

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (sc *SwapController) CreateSwapRequestHandler(c *gin.Context) {
	userID := utils.ValidateSession(c)

	var swapRequest models.SwapRequest
	utils.ParseJSONRequest(c, &swapRequest)

	// Check if there already exists a swap request for this user
	var existingRequest models.SwapRequest
	err := sc.DB.First(&existingRequest, "user_id = ?", userID).Error
	if err == nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Swap request already exists for this user")
		return
	} else if err != gorm.ErrRecordNotFound {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	// If no existing request, create a new one. Get current mess of the user
	var user models.User
	err = sc.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user info")
		return
	}

	// Get current mess of the user and set direction accordingly
	switch user.Mess {
	case 1, 2:
		swapRequest.Direction = "A to B"
	case 3, 4:
		swapRequest.Direction = "B to A"
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "User is not assigned to any mess")
		return
	}

	swapRequest.UserID = userID

	err = sc.DB.Create(&swapRequest).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create swap request")
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, models.APIResponse{
		Message: "Swap request created successfully",
		Data:    swapRequest,
	})
}
