package hosteloffice

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// BulkEditStudents handles bulk editing of multiple students
func (oc *OfficeController) BulkEditStudents(c *gin.Context) {
	var request models.BulkEditRequest

	if err := utils.ParseJSONRequest(c, &request); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	var errors []string
	recordsUpdated := 0
	updateTime := time.Now()

	// Process each update
	for _, update := range request.Updates {
		result := oc.DB.Model(&models.User{}).
			Where("roll_no = ?", update.RollNo).
			Updates(map[string]interface{}{
				"mess":         update.Mess,
				"can_register": update.CanRegister,
				"updated_at":   updateTime,
			})

		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("Failed to update user %s: %v", update.RollNo, result.Error))
		} else if result.RowsAffected == 0 {
			errors = append(errors, fmt.Sprintf("User %s not found", update.RollNo))
		} else {
			recordsUpdated++
		}
	}

	response := models.UploadResponse{
		Message:       "Bulk update processed",
		RecordsUpdate: recordsUpdated,
		Errors:        errors,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}
