package registration

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (m *MessController) VegMessRegistrationHandler(c *gin.Context) {
	// Only accept requests on this endpoint at a specified date
	// Check if the current date is within the registration period

	if !utils.GetVegRegistrationStatus(m.DB) {
		utils.RespondWithJSON(c, http.StatusForbidden, "Registration Has Ended.")
		return
	}
	userID := utils.ValidateSession(c)

	// Check if user exists and can register (from database)
	var user models.User
	if err := m.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		return
	}

	if !user.CanRegister {
		utils.RespondWithError(c, http.StatusBadRequest, "User cannot Register.")
		return
	}

	// Check if user already has a mess assigned for next period (check both DB and Redis)
	if user.NextMess != 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "User already has a mess assigned for next period")
		return
	}

	// Also check Redis for any pending assignment
	redisMessID, err := m.redisService.GetUserMess(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to check Redis assignment: "+err.Error())
		return
	}
	if redisMessID != 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "User already has a mess assignment pending")
		return
	}

	// Veg mess is always mess 5
	allowedMesses := []int{5} // Baadme logic expand kar sakte if galti se bhi koi aur mess veg option de diya toh

	var registrationErrs []string
	registered := false

	for _, mess := range allowedMesses {
		// Attempt registration using Redis (atomic operation)
		success, err := m.redisService.AttemptMessRegistration(userID, mess)
		if err != nil {
			registrationErrs = append(registrationErrs, "Mess "+strconv.Itoa(mess)+": "+err.Error())
			continue
		}

		if success {
			registered = true
			break
		} else {
			registrationErrs = append(registrationErrs, "Mess "+strconv.Itoa(mess)+": Veg Registration failed due to capacity or conflict")
		}
	}

	if !registered {
		utils.RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("Veg Registration failed: %v", registrationErrs))
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Mess registration successful. Changes will be synced to database shortly.",
	})
}
