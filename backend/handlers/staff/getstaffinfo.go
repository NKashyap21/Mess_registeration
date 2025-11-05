package staff

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (sc *ScanningController) GetStaffInfo(c *gin.Context) {
	// Get the authenticated mess staff user from context
	staffUser, exists := c.Get("user")
	if !exists {
		utils.RespondWithJSON(c, http.StatusUnauthorized, models.APIResponse{
			Message: "Authentication required",
		})
		return
	}

	staff := staffUser.(models.User)

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Staff information retrieved successfully",
		Data:    map[string]interface{}{"staff": staff},
	})
}
