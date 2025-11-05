package registration

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (m *MessController) RefreshCapacitiesHandler(c *gin.Context) {
	if err := m.redisService.RefreshCapacitiesFromDB(); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to refresh capacities: "+err.Error())
		return
	}

	// Get updated stats to return
	stats, err := m.redisService.GetAllMessStats()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated statistics: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Capacities refreshed successfully",
		Data: map[string]interface{}{
			"stats": stats,
		},
	})
}
