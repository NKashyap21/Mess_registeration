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

// DownloadRegistrationsCSV exports registrations for a date range
func (oc *OfficeController) DownloadRegistrationsCSV(c *gin.Context) {
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

	var registrations []utils.RegistrationRecord

	// Query users who registered (have next_mess set) within the date range using GORM
	var users []models.User
	if err := oc.DB.Select("id", "name", "roll_no", "email", "mess", "next_mess", "updated_at").
		Where("updated_at BETWEEN ? AND ? AND type = ?", fromDate, toDate.Add(24*time.Hour), 0).
		Order("updated_at DESC").
		Find(&users).Error; err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to fetch registrations")
		return
	}

	// Convert to RegistrationRecord format
	for _, user := range users {
		registrations = append(registrations, utils.RegistrationRecord{
			UserID:    user.ID,
			UserName:  user.Name,
			RollNo:    user.RollNo,
			Email:     user.Email,
			Mess:      user.Mess,
			NextMess:  user.NextMess,
			UpdatedAt: user.UpdatedAt,
		})
	}

	buffer, err := utils.ExportRegistrationsToCSV(registrations)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate CSV")
		return
	}

	filename := fmt.Sprintf("registrations_%s_to_%s.csv", request.FromDate, request.ToDate)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}
