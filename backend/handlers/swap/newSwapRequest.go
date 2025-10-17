package swap

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/db"
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
	case 5:
		utils.RespondWithError(c, http.StatusBadRequest, "User is in veg mess, cannot create swap request")
		return
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "User is not assigned to any mess")
		return
	}

	// Check if there is a public request already floating with the opposite direction
	// If yes, auto accept both requests and update mess of both users
	var oppositeRequest models.SwapRequest
	err = sc.DB.First(&oppositeRequest, "direction = ? AND completed = ?",
		map[string]string{"A to B": "B to A", "B to A": "A to B"}[swapRequest.Direction], false).Error
	if err == nil {
		// Found an opposite request, auto accept both
		swapRequest.Completed = true
		oppositeRequest.Completed = true
		err = sc.DB.Update("completed", oppositeRequest).Error
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update opposite swap request")
			return
		}

		// Update mess of both users
		var oppositeUser models.User
		err = sc.DB.First(&oppositeUser, "id = ?", oppositeRequest.UserID).Error
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch opposite user info")
			return
		}
		
		user.Mess, oppositeUser.Mess = oppositeUser.Mess, user.Mess
		err = sc.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(&user).Error; err != nil {
				return err
			}
			if err := tx.Save(&oppositeUser).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update users' mess")
			return
		}
	}

	var dbSwapRequest db.SwapRequest

	dbSwapRequest.Type = swapRequest.Type
	dbSwapRequest.UserID = userID
	dbSwapRequest.Direction = swapRequest.Direction
	dbSwapRequest.Completed = swapRequest.Completed
	dbSwapRequest.Password = swapRequest.Password

	err = sc.DB.Create(&dbSwapRequest).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create swap request")
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, models.APIResponse{
		Message: "Swap request created successfully",
		Data:    swapRequest,
	})
}
