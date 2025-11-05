package registration

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (m *MessController) GetMessStatsGroupedHandler(c *gin.Context) {
	stats, err := m.redisService.GetMessStatsByGroup()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get grouped mess statistics: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"stats": stats,
		},
	})
}
