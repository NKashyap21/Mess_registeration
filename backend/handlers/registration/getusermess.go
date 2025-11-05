package registration

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/services"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (m *MessController) GetUserMessHandler(c *gin.Context) {
	userID := utils.ValidateSession(c)

	// First check Redis for any pending assignment
	redisMessID, err := m.redisService.GetUserMess(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to check Redis assignment: "+err.Error())
		return
	}

	if redisMessID != 0 {
		utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
			Data: map[string]interface{}{
				"mess":   redisMessID,
				"status": "pending_sync",
			},
		})
		return
	}

	// Check database
	var user models.User
	if err := m.DB.First(&user, userID).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch user: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"current_mess":      user.Mess,
			"current_mess_name": services.GetMessName(int(user.Mess)),
			"next_mess":         user.NextMess,
			"next_mess_name":    services.GetMessName(int(user.NextMess)),
			"status":            "confirmed",
		},
	})
}
