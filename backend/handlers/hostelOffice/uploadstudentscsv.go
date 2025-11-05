package hosteloffice

import (
	"fmt"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// UploadStudentsCSV handles bulk upload of students via CSV file
// Expected CSV format: Name,Email,Phone,RollNo,Mess,Type,CanRegister
func (oc *OfficeController) UploadStudentsCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	users, err := utils.ParseStudentsCSV(file)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("Error parsing CSV: %v", err))
		return
	}

	var errors []string
	recordsAdded := 0

	// Process each user
	for _, user := range users {
		if err := oc.DB.Create(&user).Error; err != nil {
			errors = append(errors, fmt.Sprintf("Failed to add user %s: %v", user.RollNo, err))
		} else {
			recordsAdded++
		}
	}

	response := models.UploadResponse{
		Message:      "CSV upload processed",
		RecordsAdded: recordsAdded,
		Errors:       errors,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}
