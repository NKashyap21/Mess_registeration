package hosteloffice

import (
	"log"
	"net/http"
	"time"

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
	var reg models.MessRegistrationDetails

	// There should only be one row
	if err := oc.DB.First(&reg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new record with NormalRegistrationOpen = true by default
			reg = models.MessRegistrationDetails{NormalRegistrationOpen: true}
			if err := oc.DB.Create(&reg).Error; err != nil {
				utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create registration details")
				return
			}
		} else {
			log.Println("Error fetching registration details:", err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
			return
		}
	} else {
		// Flip the normal registration state
		reg.NormalRegistrationOpen = !reg.NormalRegistrationOpen
		if err := oc.DB.Save(&reg).Error; err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to toggle normal registration")
			return
		}
	}

	status := "closed"
	if reg.NormalRegistrationOpen {
		status = "open"
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message":                  "Normal registration toggled successfully",
		"status":                   status,
		"normal_registration_open": reg.NormalRegistrationOpen,
		"veg_registration_open":    reg.VegRegistrationOpen,
		"mess_a_ldh_capacity":      reg.MessALDHCapacity,
		"mess_a_udh_capacity":      reg.MessAUDHCapacity,
		"mess_b_ldh_capacity":      reg.MessBLDHCapacity,
		"mess_b_udh_capacity":      reg.MessBUDHCapacity,
		"veg_mess_capacity":        reg.VegMessCapacity,
	})
}

func (oc *OfficeController) ToggleVegRegistration(c *gin.Context) {
	var reg models.MessRegistrationDetails

	if err := oc.DB.First(&reg).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new record with VegRegistrationOpen = true by default
			reg = models.MessRegistrationDetails{VegRegistrationOpen: true}
			if err := oc.DB.Create(&reg).Error; err != nil {
				utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create registration details")
				return
			}
		} else {
			log.Println("Error fetching registration details:", err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
			return
		}
	} else {
		// Flip the veg registration state
		reg.VegRegistrationOpen = !reg.VegRegistrationOpen
		if err := oc.DB.Save(&reg).Error; err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to toggle veg registration")
			return
		}
	}

	status := "closed"
	if reg.VegRegistrationOpen {
		status = "open"
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"message":                  "Veg registration toggled successfully",
		"status":                   status,
		"normal_registration_open": reg.NormalRegistrationOpen,
		"veg_registration_open":    reg.VegRegistrationOpen,
		"mess_a_ldh_capacity":      reg.MessALDHCapacity,
		"mess_a_udh_capacity":      reg.MessAUDHCapacity,
		"mess_b_ldh_capacity":      reg.MessBLDHCapacity,
		"mess_b_udh_capacity":      reg.MessBUDHCapacity,
		"veg_mess_capacity":        reg.VegMessCapacity,
	})
}

func (oc *OfficeController) GetRegistrationStatus(c *gin.Context) {
	var details models.MessRegistrationDetails

	// Fetch the single record (you can use LIMIT 1 or WHERE 1=1)
	if err := oc.DB.First(&details).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not fetch registration status")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"veg_registration_open":    details.VegRegistrationOpen,
		"normal_registration_open": details.NormalRegistrationOpen,
		"mess_a_ldh_capacity":      details.MessALDHCapacity,
		"mess_a_udh_capacity":      details.MessAUDHCapacity,
		"mess_b_ldh_capacity":      details.MessBLDHCapacity,
		"mess_b_udh_capacity":      details.MessBUDHCapacity,
		"veg_mess_capacity":        details.VegMessCapacity,
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
