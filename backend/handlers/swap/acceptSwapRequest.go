package swap

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (sc *SwapController) AcceptSwapRequestHandler(c *gin.Context) {
	userID := utils.ValidateSession(c)

	var selectedSwapRequest db.SwapRequest
	utils.ParseJSONRequest(c, &selectedSwapRequest)

	var user models.User
	if err := sc.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "User not found")
		return
	}
	
	// User does not have a mess assigned, cannot swap to null
	if user.Mess == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "User does not have a mess assigned")
		return
	}
	
	// Check if user has already swapped mess once
	var existingSwap db.SwapRequest
	if err := sc.DB.First(&existingSwap, "user_id = ? AND completed = ?", userID, true).Error; err == nil {
		utils.RespondWithError(c, http.StatusForbidden, "User has already swapped mess once")
		return
	}

	// Two options, either the user is accepting a friend request or a public request
	switch selectedSwapRequest.Type {
	case "friend":
		sc.acceptFriendSwapRequest(userID, selectedSwapRequest, c)
	case "public":
		sc.acceptPublicSwapRequest(userID, selectedSwapRequest, c)
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid swap request type")
		return
	}
}

func (sc *SwapController) acceptPublicSwapRequest(userID uint, selectedSwapRequest db.SwapRequest, c *gin.Context) {
	// Check if the selected swap request exists and is of type 'public'
	var existingRequest db.SwapRequest
	if err := sc.DB.First(&existingRequest, "user_id = ? AND type = ?", selectedSwapRequest.UserID, "public").Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Swap request not found")
		return
	}

	// Check if the user is trying to accept their own request
	if existingRequest.UserID == userID {
		utils.RespondWithError(c, http.StatusForbidden, "Cannot accept your own swap request")
		return
	}

	// Check if user has the same mess as the existing request
	var user models.User
	if err := sc.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "User not found")
		return
	}

	var theoreticalSwapDirection string
	switch user.Mess {
	case 1, 2:
		theoreticalSwapDirection = "A to B"
	case 3, 4:
		theoreticalSwapDirection = "B to A"
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "User does not belong to a valid mess")
		return
	}

	if theoreticalSwapDirection == existingRequest.Direction {
		utils.RespondWithError(c, http.StatusBadRequest, "Cannot swap with the same mess")
		return
	}

	// Perform the swap
	if err := sc.performSwap(userID, existingRequest.UserID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to perform swap")
		return
	}

	// Delete the accepted swap request
	if err := sc.DB.Where("user_id = ?", existingRequest.UserID).Delete(&existingRequest).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete swap request")
		return
	}

	utils.RespondWithJSON(c, http.StatusAccepted, models.APIResponse{
		Message: "Public swap request accepted and performed successfully",
	})
}

func (sc *SwapController) acceptFriendSwapRequest(userID uint, selectedSwapRequest db.SwapRequest, c *gin.Context) {
	// Check if the selected swap request exists and is of type 'friend'
	var existingRequest db.SwapRequest
	if err := sc.DB.First(&existingRequest, "user_id = ? AND type = ?", selectedSwapRequest.UserID, "friend").Error; err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Swap request not found")
		return
	}

	// Check if the user is trying to accept their own request
	if existingRequest.UserID == userID {
		utils.RespondWithError(c, http.StatusForbidden, "Cannot accept your own swap request")
		return
	}

	// Check if user has the same mess as the existing request
	var user models.User
	if err := sc.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "User not found")
		return
	}

	var theoreticalSwapDirection string
	switch user.Mess {
	case 1, 2:
		theoreticalSwapDirection = "A to B"
	case 3, 4:
		theoreticalSwapDirection = "B to A"
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "User does not belong to a valid mess")
		return
	}

	if theoreticalSwapDirection == existingRequest.Direction {
		utils.RespondWithError(c, http.StatusBadRequest, "Cannot swap with the same mess")
		return
	}

	// Verify the password
	if existingRequest.Password != selectedSwapRequest.Password {
		utils.RespondWithError(c, http.StatusForbidden, "Incorrect password for friend swap request")
		return
	}

	// Perform the swap
	if err := sc.performSwap(userID, existingRequest.UserID); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to perform swap")
		return
	}

	// Delete the accepted swap request
	if err := sc.DB.Where("user_id = ?", existingRequest.UserID).Delete(&existingRequest).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete swap request")
		return
	}

	utils.RespondWithJSON(c, http.StatusAccepted, models.APIResponse{
		Message: "Friend swap request accepted and performed successfully",
	})
}

func (sc *SwapController) performSwap(userID1, userID2 uint) error {
	var user1, user2 models.User
	if err := sc.DB.First(&user1, userID1).Error; err != nil {
		return err
	}
	if err := sc.DB.First(&user2, userID2).Error; err != nil {
		return err
	}

	// Swap the mess assignments
	user1.Mess, user2.Mess = user2.Mess, user1.Mess

	if err := sc.DB.Save(&user1).Error; err != nil {
		return err
	}
	if err := sc.DB.Save(&user2).Error; err != nil {
		return err
	}

	return nil
}
