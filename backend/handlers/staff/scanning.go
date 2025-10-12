package staff

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		Data:    gin.H{"staff": staff},
	})
}

func (sc *ScanningController) ScanningHandler(c *gin.Context) {
	// Get the authenticated mess staff user from context
	staffUser, exists := c.Get("user")
	if !exists {
		utils.RespondWithJSON(c, http.StatusUnauthorized, models.APIResponse{
			Message: "Authentication required",
		})
		return
	}

	staff := staffUser.(models.User)

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
	// The mess staff should only be able to scan users from their assigned mess
	if staff.Mess == 0 {
		utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
			Message: "Staff does not have a mess assigned",
		})
		return
	}

	// Check mess access based on staff's assigned mess
	switch staff.Mess {
	case 1, 2: // Mess A (LDH & UDH)
		if user.Mess != 1 && user.Mess != 2 {
			utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
				Message: "User does not have access to Mess A",
			})
			return
		}
	case 3, 4: // Mess B (LDH & UDH)
		if user.Mess != 3 && user.Mess != 4 {
			utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
				Message: "User does not have access to Mess B",
			})
			return
		}
	default:
		utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
			Message: "Invalid mess assignment for staff",
		})
		return
	}

	// If all checks pass, respond with user details
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "User verified successfully",
		Data:    gin.H{"user": user, "staff": staff.Name},
	})
}
