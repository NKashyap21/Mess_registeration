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

// DownloadScansCSV exports scans for a date range
func (oc *OfficeController) DownloadScansCSV(c *gin.Context) {
	var request models.DateRangeRequest

	if err := c.ShouldBindQuery(&request); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid date range parameters")
		return
	}

	fromDate, err := time.Parse("2006-01-02", request.FromDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid from_date format. Use YYYY-MM-DD")
		return
	}

	toDate, err := time.Parse("2006-01-02", request.ToDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid to_date format. Use YYYY-MM-DD")
		return
	}

	var scans []utils.ScanRecord

	// Query scans with user information using GORM joins
	var scanResults []struct {
		ID        uint      `gorm:"column:id"`
		UserID    uint      `gorm:"column:user_id"`
		UserName  string    `gorm:"column:user_name"`
		RollNo    string    `gorm:"column:roll_no"`
		MessID    uint      `gorm:"column:mess_id"`
		Meal      int       `gorm:"column:meal"`
		Date      time.Time `gorm:"column:date"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	if err := oc.DB.Table("scans s").
		Select("s.id, s.user_id, u.name as user_name, u.roll_no, s.mess_id, s.meal, s.date, s.created_at").
		Joins("JOIN users u ON s.user_id = u.id").
		Where("s.date BETWEEN ? AND ?", fromDate, toDate).
		Order("s.date DESC, s.created_at DESC").
		Find(&scanResults).Error; err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch scans")
		return
	}

	// Convert to ScanRecord format
	for _, sr := range scanResults {
		scans = append(scans, utils.ScanRecord{
			ID:        sr.ID,
			UserID:    sr.UserID,
			UserName:  sr.UserName,
			RollNo:    sr.RollNo,
			MessID:    sr.MessID,
			Meal:      sr.Meal,
			Date:      sr.Date,
			CreatedAt: sr.CreatedAt,
		})
	}

	buffer, err := utils.ExportScansToCSV(scans)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	filename := fmt.Sprintf("scans_%s_to_%s.csv", request.FromDate, request.ToDate)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}
