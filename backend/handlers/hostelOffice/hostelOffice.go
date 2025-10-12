package hosteloffice

import (
	"log"
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OfficeController struct {
	DB *gorm.DB
}

type UserInfo struct {
	Name   string
	Email  string
	RollNo string
	Mess   int8
}

type EditUserInfo struct {
	RollNo      string `json:"roll_no"`      //To identify the student
	Mess        int8   `json:"mess"`         //The hostel office can change this value
	CanRegister bool   `json:"can_register"` //fasle -> The user has been deactivated.
}

func InitOfficeController() *OfficeController {
	return &OfficeController{
		DB: config.GetDB(),
	}
}

func (oc *OfficeController) GetStudents(c *gin.Context) {
	//Returns the name,roll_no,email,mess of all the studnets.
	var studentsInfo []UserInfo

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

	var stundentInfo UserInfo
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
	var studentInfo EditUserInfo

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

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Success Updating the usef information."})

}
