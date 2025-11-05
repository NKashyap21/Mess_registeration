package hosteloffice

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// DownloadStudentsCSV exports all students to CSV
func (oc *OfficeController) DownloadStudentsCSV(c *gin.Context) {
	var users []models.User

	if err := oc.DB.Find(&users).Error; err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	buffer, err := utils.ExportUsersToCSV(users)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=students_%s.csv", time.Now().Format("2006-01-02")))
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}
