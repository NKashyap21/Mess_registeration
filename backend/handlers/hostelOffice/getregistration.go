package hosteloffice

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (oc *OfficeController) GetRegistrationStatus(c *gin.Context) {
	var details models.MessRegistrationDetails

	if err := oc.DB.First(&details).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not fetch registration details")
		return
	}

	type StatsResult struct {
		Mess          int8
		CurrentCount  int64
		UpcomingCount int64
	}

	currentStats := map[int]int{0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	upcomingStats := map[int]int{0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0}

	// Count current mess assignments for each mess (1-5)
	for mess := 1; mess <= 5; mess++ {
		var count int64
		oc.DB.Model(&models.User{}).
			Where("mess = ? AND can_register = ?", mess, true).
			Count(&count)
		currentStats[mess] = int(count)
	}

	// Count upcoming mess assignments for each mess (1-5)
	for mess := 1; mess <= 5; mess++ {
		var count int64
		oc.DB.Model(&models.User{}).
			Where("next_mess = ? AND can_register = ?", mess, true).
			Count(&count)
		upcomingStats[mess] = int(count)
	}

	// Compute unassigned_capacity, total_can_register, total_registered_yet
	var totalCanRegister int64
	var totalStudents int64
	var totalCanButDidntRegister int64
	var totalCanButDidntRegisterUpcoming int64

	oc.DB.Model(&models.User{}).Where("type = ?", 0).Count(&totalStudents) // total students
	oc.DB.Model(&models.User{}).Where("can_register = ?", true).Count(&totalCanRegister)
	oc.DB.Model(&models.User{}).Where("mess = 0 AND can_register = true").Count(&totalCanButDidntRegister)
	oc.DB.Model(&models.User{}).Where("next_mess = 0 AND can_register = true").Count(&totalCanButDidntRegisterUpcoming)

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Data: map[string]interface{}{
			"registration_status": map[string]interface{}{
				"veg":    details.VegRegistrationOpen,
				"normal": details.NormalRegistrationOpen,
			},
			"capacity": map[string]interface{}{
				"Mess A": map[string]interface{}{
					"LDH": details.MessALDHCapacity,
					"UDH": details.MessAUDHCapacity,
				},
				"Mess B": map[string]interface{}{
					"LDH": details.MessBLDHCapacity,
					"UDH": details.MessBUDHCapacity,
				},
				"Extra": map[string]interface{}{
					"Veg Mess":   details.VegMessCapacity,
					"Unassigned": totalCanRegister,
				},
			},
			"current_mess": map[string]interface{}{
				"Mess A": map[string]interface{}{
					"LDH": currentStats[1],
					"UDH": currentStats[2],
				},
				"Mess B": map[string]interface{}{
					"LDH": currentStats[3],
					"UDH": currentStats[4],
				},
				"Extra": map[string]interface{}{
					"Veg Mess":   currentStats[5],
					"Unassigned": totalCanButDidntRegister,
				},
			},
			"upcoming_mess": map[string]interface{}{
				"Mess A": map[string]interface{}{
					"LDH": upcomingStats[1],
					"UDH": upcomingStats[2],
				},
				"Mess B": map[string]interface{}{
					"LDH": upcomingStats[3],
					"UDH": upcomingStats[4],
				},
				"Extra": map[string]interface{}{
					"Veg Mess":   upcomingStats[5],
					"Unassigned": totalCanButDidntRegisterUpcoming,
				},
			},
		},
	})

}
