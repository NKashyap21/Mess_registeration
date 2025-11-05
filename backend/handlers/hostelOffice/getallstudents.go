package hosteloffice

import (
	"log"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
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
		utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{Message: "No students found."})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, studentsInfo)
}