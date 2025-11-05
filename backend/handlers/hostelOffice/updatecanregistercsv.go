package hosteloffice

import (
	"fmt"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// UpdateCanRegisterCSV handles bulk update of can_register field via CSV file
// Expected CSV format: RollNo,CanRegister
func (oc *OfficeController) UpdateCanRegisterCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	updates, err := utils.ParseCanRegisterCSV(file)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, fmt.Sprintf("Error parsing CSV: %v", err))
		return
	}

	var errors []string
	recordsUpdated := 0

	// Process each update
	for rollNo, canRegister := range updates {
		result := oc.DB.Model(&models.User{}).
			Where("roll_no = ?", rollNo).
			Update("can_register", canRegister)

		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("Failed to update user %s: %v", rollNo, result.Error))
		} else if result.RowsAffected == 0 {
			errors = append(errors, fmt.Sprintf("User %s not found", rollNo))
		} else {
			recordsUpdated++
		}
	}

	response := models.UploadResponse{
		Message:       "CSV upload processed",
		RecordsUpdate: recordsUpdated,
		Errors:        errors,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}
