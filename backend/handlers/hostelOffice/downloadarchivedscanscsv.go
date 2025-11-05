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

// DownloadArchivedScansCSV exports data from an archived scans table
func (oc *OfficeController) DownloadArchivedScansCSV(c *gin.Context) {
	tableName := c.Query("table")
	if tableName == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "table query parameter is required")
		return
	}

	var scans []utils.ScanRecord

	// Query archived scans using GORM Table() method
	var scanResults []struct {
		ID        uint      `gorm:"column:id"`
		UserID    uint      `gorm:"column:user_id"`
		MessID    uint      `gorm:"column:mess_id"`
		Meal      int       `gorm:"column:meal"`
		Date      time.Time `gorm:"column:date"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	if err := oc.DB.Table(tableName).
		Select("id", "user_id", "mess_id", "meal", "date", "created_at").
		Order("date DESC, created_at DESC").
		Find(&scanResults).Error; err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch archived scans")
		return
	}

	// Convert to ScanRecord and enrich with user data
	for _, sr := range scanResults {
		scanRecord := utils.ScanRecord{
			ID:        sr.ID,
			UserID:    sr.UserID,
			UserName:  "",
			RollNo:    "",
			MessID:    sr.MessID,
			Meal:      sr.Meal,
			Date:      sr.Date,
			CreatedAt: sr.CreatedAt,
		}

		// Try to enrich with user data from current users table
		var user models.User
		if err := oc.DB.Select("name", "roll_no").Where("id = ?", sr.UserID).First(&user).Error; err == nil {
			scanRecord.UserName = user.Name
			scanRecord.RollNo = user.RollNo
		}

		scans = append(scans, scanRecord)
	}

	buffer, err := utils.ExportScansToCSV(scans)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", tableName))
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}
