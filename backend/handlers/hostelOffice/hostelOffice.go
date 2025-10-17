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

func (oc *OfficeController) StartRegRegistration(c *gin.Context) {
	var registrationDetails models.MessRegistrationDetails

	if err := utils.ParseJSONRequest(c, &registrationDetails); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	// Check if payload tries to change veg registration dates
	if !registrationDetails.VegRegistrationStart.IsZero() || !registrationDetails.VegRegistrationEnd.IsZero() {
		utils.RespondWithError(c, http.StatusBadRequest, "Cannot change veg registration dates in normal registration endpoint")
		return
	}

	// Update the dates in DB
	if err := oc.DB.First(&registrationDetails, "WHERE 1=1").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new record
			if err := oc.DB.Create(&registrationDetails).Error; err != nil {
				log.Println(err)
				utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		} else {
			log.Println(err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
}

func (oc *OfficeController) StartVegRegistration(c *gin.Context) {
	var registrationDetails models.MessRegistrationDetails

	if err := utils.ParseJSONRequest(c, &registrationDetails); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	// Check if payload tries to change normal registration dates
	if !registrationDetails.NormalRegistrationStart.IsZero() || !registrationDetails.NormalRegistrationEnd.IsZero() {
		utils.RespondWithError(c, http.StatusBadRequest, "Cannot change normal registration dates in veg registration endpoint")
		return
	}

	// Update the dates in DB
	if err := oc.DB.First(&registrationDetails, "WHERE 1=1").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new record
			if err := oc.DB.Create(&registrationDetails).Error; err != nil {
				log.Println(err)
				utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
		} else {
			log.Println(err)
			utils.RespondWithError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}
}