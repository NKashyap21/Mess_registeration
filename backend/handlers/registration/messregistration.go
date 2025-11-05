package registration

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/services"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (m *MessController) MessRegistrationHandler(c *gin.Context) {
	logger := services.GetLoggerService()

	if !utils.GetNormalRegistrationStatus(m.DB) {
		utils.RespondWithError(c, http.StatusForbidden, "Registration Has Ended.")
		return
	}

	userID := utils.ValidateSession(c)

	// Check if user exists and can register (from database)
	var user models.User
	if err := m.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		logger.LogDatabaseAction(userID, "SELECT", "users", fmt.Sprintf("Failed to fetch user: %s", err.Error()), c.ClientIP())
		return
	}

	if !user.CanRegister {
		utils.RespondWithError(c, http.StatusBadRequest, "User cannot Register.")
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", "User cannot register (CanRegister=false)", c.ClientIP())
		return
	}

	// Check if user already has a mess assigned for next period (check both DB and Redis)
	if user.NextMess != 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "User already has a mess assigned for next period")
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("User already has NextMess %d assigned", user.NextMess), c.ClientIP())
		return
	}

	// Also check Redis for any pending assignment
	redisMessID, err := m.redisService.GetUserMess(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to check Redis assignment: "+err.Error())
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("Redis check failed: %s", err.Error()), c.ClientIP())
		return
	}
	if redisMessID != 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "User already has a mess assignment pending")
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("User already has pending assignment for mess %d", redisMessID), c.ClientIP())
		return
	}

	messParam := c.Param("mess")
	mess, err := strconv.Atoi(messParam)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Mess parameter must be an integer")
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("Invalid mess parameter: %s", messParam), c.ClientIP())
		return
	}

	if !services.IsValidMessID(mess) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid mess ID. Valid options: 1 (MessA LDH), 2 (MessA UDH), 3 (MessB LDH), 4 (MessB UDH)")
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("Invalid mess ID: %d", mess), c.ClientIP())
		return
	}

	// Attempt registration using Redis (atomic operation)
	success, err := m.redisService.AttemptMessRegistration(userID, mess)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Registration failed: "+err.Error())
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("Redis registration failed for mess %d: %s", mess, err.Error()), c.ClientIP())
		return
	}

	if !success {
		utils.RespondWithError(c, http.StatusBadRequest, "Registration failed due to capacity or conflict")
		logger.LogUserAction(userID, "MESS_REGISTRATION_FAILED", fmt.Sprintf("Registration failed for mess %d due to capacity or conflict", mess), c.ClientIP())
		return
	}

	// Log successful registration
	logger.LogUserAction(userID, "MESS_REGISTRATION_SUCCESS", fmt.Sprintf("User successfully registered for mess %d", mess), c.ClientIP())
	logger.LogDatabaseAction(userID, "UPDATE", "mess_assignments", fmt.Sprintf("User assigned to mess %d via Redis", mess), c.ClientIP())

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Mess registration successful. Changes will be synced to database shortly.",
	})
}
