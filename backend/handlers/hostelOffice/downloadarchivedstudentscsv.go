package hosteloffice

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// DownloadArchivedStudentsCSV exports data from an archived users table
func (oc *OfficeController) DownloadArchivedStudentsCSV(c *gin.Context) {
	tableName := c.Query("table")
	if tableName == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "table query parameter is required")
		return
	}

	var users []models.User
	if err := utils.GetArchivedTableData(oc.DB, tableName, &users); err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch archived data")
		return
	}

	buffer, err := utils.ExportUsersToCSV(users)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", tableName))
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}
