package hosteloffice

import (
	"log"
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/db"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (oc *OfficeController) GetStudents(c *gin.Context) {
	//Returns the name,roll_no,email,mess of all the studnets.
	var studentsInfo []models.UserInfo

	result := oc.DB.Select("name", "roll_no", "email", "mess").Find(&studentsInfo)

	if result.Error != nil {
		log.Println(result.Error)
		utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//In case the students are not found.
	if result.RowsAffected == 0 {
		utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "No students found."})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, studentsInfo)
}

func (oc *OfficeController) GetStudentsByID(c *gin.Context) {
	//Returns the name,email,roll_no,mess of the student.

	roll_no := c.Param("roll_no")

	var stundentInfo models.UserInfo
	result := oc.DB.Select("name", "email", "roll_no", "mess").Where("roll_no = ?", roll_no).Find(&stundentInfo)

	if result.Error != nil {
		log.Println(result.Error)
		utils.RespondWithError(c, http.StatusInternalServerError, "Internal Status Error")
		return
	}

	if result.RowsAffected == 0 {
		utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Student not found"})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, stundentInfo)

}

func (oc *OfficeController) EditStudentById(c *gin.Context) {
	// The hostel office can change the mess and also the deactivate a student.
	var studentInfo models.EditUserInfo

	if err := utils.ParseJSONRequest(c, &studentInfo); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	updateTime := time.Now()

	result := oc.DB.Model(&models.User{}).
		Where("roll_no = ?", studentInfo.RollNo).
		Updates(map[string]interface{}{
			"mess":         studentInfo.Mess,
			"can_register": studentInfo.CanRegister,
			"updated_at":   updateTime,
		})

	if result.Error != nil {
		log.Println(result.Error)
		utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if result.RowsAffected == 0 {
		utils.RespondWithJSON(c, http.StatusBadRequest, gin.H{"message": "Student not found!"})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Success Updating the user information."})

}

// func (oc *OfficeController) StartRegRegistration(c *gin.Context) {
// 	var registrationDetails models.MessRegistrationDetails

// 	if err := utils.ParseJSONRequest(c, &registrationDetails); err != nil {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
// 		return
// 	}

// 	// Check if payload tries to change veg registration dates
// 	if !registrationDetails.VegRegistrationStart.IsZero() || !registrationDetails.VegRegistrationEnd.IsZero() {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Cannot change veg registration dates in normal registration endpoint")
// 		return
// 	}

// 	// Update the dates in DB
// 	if err := oc.DB.First(&registrationDetails, "WHERE 1=1").Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			// Create new record
// 			if err := oc.DB.Create(&registrationDetails).Error; err != nil {
// 				log.Println(err)
// 				utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
// 				return
// 			}
// 		} else {
// 			log.Println(err)
// 			utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
// 			return
// 		}
// 	}
// }

// func (oc *OfficeController) StartVegRegistration(c *gin.Context) {
// 	var registrationDetails models.MessRegistrationDetails

// 	if err := utils.ParseJSONRequest(c, &registrationDetails); err != nil {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
// 		return
// 	}

// 	// Check if payload tries to change normal registration dates
// 	if !registrationDetails.NormalRegistrationStart.IsZero() || !registrationDetails.NormalRegistrationEnd.IsZero() {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Cannot change normal registration dates in veg registration endpoint")
// 		return
// 	}

// 	// Update the dates in DB
// 	if err := oc.DB.First(&registrationDetails, "WHERE 1=1").Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			// Create new record
// 			if err := oc.DB.Create(&registrationDetails).Error; err != nil {
// 				log.Println(err)
// 				utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
// 				return
// 			}
// 		} else {
// 			log.Println(err)
// 			utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
// 			return
// 		}
// 	}
// }

func (oc *OfficeController) ToggleNormalRegistration(c *gin.Context) {
	// Toggle NormalRegistrationOpen for the first row
	if err := oc.DB.Model(&db.MessRegistrationDetails{}).
		Limit(1).
		Update("normal_registration_open", gorm.Expr("NOT normal_registration_open")).
		Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to toggle normal registration")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message": "Normal registration toggled successfully",
	})
}

func (oc *OfficeController) ToggleVegRegistration(c *gin.Context) {
	// Toggle VegRegistrationOpen for the first row
	if err := oc.DB.Model(&db.MessRegistrationDetails{}).
		Limit(1).
		Update("veg_registration_open", gorm.Expr("NOT veg_registration_open")).
		Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to toggle veg registration")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message": "Veg registration toggled successfully",
	})
}

func (oc *OfficeController) ApplyNewRegistration(c *gin.Context) {
	// Step 1: Set both registrations to false
	if err := oc.DB.Model(&models.MessRegistrationDetails{}).
		Updates(map[string]interface{}{
			"veg_registration_open":    false,
			"normal_registration_open": false,
		}).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to close registrations")
		return
	}

	// Step 2: Begin a transaction for swapping mess info
	tx := oc.DB.Begin()
	if tx.Error != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	// Step 3: Copy next_mess → mess
	if err := tx.Exec(`UPDATE users SET mess = next_mess`).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to copy next_mess to mess")
		return
	}

	// Step 4: Reset next_mess to NULL or 0
	if err := tx.Exec(`UPDATE users SET next_mess = 0`).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to reset next_mess")
		return
	}

	// Step 5: Commit transaction
	if err := tx.Commit().Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Step 6: Return new status
	var details models.MessRegistrationDetails
	if err := oc.DB.First(&details).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch updated registration status")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message":                  "New registration cycle applied successfully",
		"veg_registration_open":    details.VegRegistrationOpen,
		"normal_registration_open": details.NormalRegistrationOpen,
	})
}

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

	var stats []StatsResult

	// Single query: count current (mess) and upcoming (next_mess)
	if err := oc.DB.Model(&models.User{}).
		Select(`mess,
			COUNT(*) as current_count,
			SUM(CASE WHEN next_mess = mess THEN 1 ELSE 0 END) as upcoming_count`).
		Group("mess").
		Scan(&stats).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to compute mess stats")
		return
	}

	currentStats := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	upcomingStats := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}

	for _, s := range stats {
		currentStats[int(s.Mess)] = int(s.CurrentCount)
		upcomingStats[int(s.Mess)] = int(s.UpcomingCount)
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

	response := gin.H{
		"registration_status": gin.H{
			"veg":    details.VegRegistrationOpen,
			"normal": details.NormalRegistrationOpen,
		},
		"capacity": gin.H{
			"Mess A": gin.H{
				"LDH": details.MessALDHCapacity,
				"UDH": details.MessAUDHCapacity,
			},
			"Mess B": gin.H{
				"LDH": details.MessBLDHCapacity,
				"UDH": details.MessBUDHCapacity,
			},
			"Extra": gin.H{
				"Veg Mess":   details.VegMessCapacity,
				"Unassigned": totalCanRegister,
			},
		},
		"current_mess": gin.H{
			"Mess A": gin.H{
				"LDH": currentStats[1],
				"UDH": currentStats[2],
			},
			"Mess B": gin.H{
				"LDH": currentStats[3],
				"UDH": currentStats[4],
			},
			"Extra": gin.H{
				"Veg Mess":   currentStats[5],
				"Unassigned": totalCanButDidntRegister,
			},
		},
		"upcoming_mess": gin.H{
			"Mess A": gin.H{
				"LDH": upcomingStats[1],
				"UDH": upcomingStats[2],
			},
			"Mess B": gin.H{
				"LDH": upcomingStats[3],
				"UDH": upcomingStats[4],
			},
			"Extra": gin.H{
				"Veg Mess":   upcomingStats[5],
				"Unassigned": totalCanButDidntRegisterUpcoming,
			},
		},
	}

	utils.RespondWithJSON(c, http.StatusOK, response)

}
