package staff

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/config"
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
			Data:    gin.H{"user": user},
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
		allowedMesses := []int8{1, 2, 5}
		if !utils.Contains(allowedMesses, user.Mess) {
			utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
				Message: "User does not have access to Mess A",
				Data:    gin.H{"user": user},
			})
			return
		}
	case 3, 4: // Mess B (LDH & UDH)
		allowedMesses := []int8{3, 4}
		if !utils.Contains(allowedMesses, user.Mess) {
			utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
				Message: "User does not have access to Mess B",
				Data:    gin.H{"user": user},
			})
			return
		}
	default:
		utils.RespondWithJSON(c, http.StatusForbidden, models.APIResponse{
			Message: "Invalid mess assignment for staff",
		})
		return
	}

	// Check if user has already scanned (check Redis)
	ctx := context.Background()
	redisClient := config.GetRedisClient()
	scanKey := fmt.Sprintf("scan:%s", rollNo)

	// Check if user has already scanned for this meal
	existsCount, err := redisClient.Exists(ctx, scanKey).Result()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Redis error: "+err.Error())
		return
	}

	if existsCount > 0 {
		utils.RespondWithJSON(c, http.StatusConflict, models.APIResponse{
			Message: "User has already scanned their ID card",
			Data:    gin.H{"user": user, "already_scanned": true},
		})
		return
	}

	// Determine the TTL based on the current time
	istLocation := time.FixedZone("IST", 5*60*60+30*60)
	currentTime := time.Now().In(istLocation)

	var ttl time.Duration
	if currentTime.Hour() == 17 { // 5-6 PM IST (snacks time)
		ttl = 1 * time.Hour
	} else {
		ttl = 150 * time.Minute
	}

	// Store the scan in Redis with the determined TTL
	scanData := fmt.Sprintf("scanned_by:%s_at:%d", staff.Name, time.Now().Unix())

	err = redisClient.Set(ctx, scanKey, scanData, ttl).Err()
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to store scan record: "+err.Error())
		return
	}

	// If all checks pass and scan is stored successfully, respond with user details
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "User verified and scan recorded successfully",
		Data:    gin.H{"user": user, "staff": staff.Name, "scan_recorded": true},
	})
}
