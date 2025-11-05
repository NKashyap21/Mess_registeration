package hosteloffice

import (
	"log"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)



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
		utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{Message: "Student not found"})
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, stundentInfo)

}