package staff

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScanningController struct {
	DB *gorm.DB
}

func InitStaffController() *ScanningController {
	return &ScanningController{
		DB: config.GetDB(),
	}
}

func (sc *ScanningController) VerifyAPIKey(c *gin.Context) string {
	// Check API KEY
	apiKey := c.GetHeader("X-API-KEY")
	if apiKey == "" {
		utils.RespondWithJSON(c, http.StatusUnauthorized, models.APIResponse{
			Message: "No API key provided",
		})
		return ""
	}
	return apiKey
}

func (sc *ScanningController) ScanningHandler(c *gin.Context) {
	// Verify Mess Staff authentication somehow
	apiKey := sc.VerifyAPIKey(c)

	// Get roll number from query parameters
	rollNo := c.Query("roll_no")
	if rollNo == "" {
		utils.RespondWithJSON(c, http.StatusBadRequest, models.APIResponse{
			Message: "No roll number entered",
		})
		return
	}

	// Fetch user details from the database
	var user models.User
	if err := sc.DB.Where("roll_no = ?", rollNo).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondWithJSON(c, http.StatusNotFound, models.APIResponse{
				Message: "User not found",
			})
		} else {
			utils.RespondWithError(c, http.StatusInternalServerError, "Database error: "+err.Error())
		}
		return
	}

	// Check if user has the correct mess assigned
	if user.Mess == 0 {
		utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
			Message: "User does not have a mess assigned",
		})
		return
	}

	// Check if scanned user has the correct mess assigned
	// For messStaff user, we can get the mess from the API KEY in the request
	if apiKey == config.GetAPIKeys().MessA {
		if user.Mess != 1 && user.Mess != 2 {
			utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
				Message: "User does not have access to Mess A",
			})
			return
		}
	} else if apiKey == config.GetAPIKeys().MessB {
		if user.Mess != 3 && user.Mess != 4 {
			utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
				Message: "User does not have access to Mess B",
			})
			return
		}
	} else {
		utils.RespondWithJSON(c, http.StatusUnauthorized, models.APIResponse{
			Message: "Invalid API key",
		})
		return
	}

	// If all checks pass, respond with user details
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "User verified successfully",
	})

}
