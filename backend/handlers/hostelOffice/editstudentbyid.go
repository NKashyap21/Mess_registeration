package hosteloffice

import (
	"log"
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

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
		utils.RespondWithJSON(c, http.StatusBadRequest, models.APIResponse{Message: "Student not found!"})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{Message: "Success Updating the user information."})

}
