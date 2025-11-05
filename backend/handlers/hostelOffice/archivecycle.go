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

// ArchiveCycle archives the current cycle data and prepares for a new cycle
func (oc *OfficeController) ArchiveCycle(c *gin.Context) {
	var request models.ArchiveRequest

	if err := utils.ParseJSONRequest(c, &request); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid payload")
		return
	}

	// Create archive date from month and year
	archiveDate := time.Date(request.Year, time.Month(request.Month), 1, 0, 0, 0, 0, time.UTC)

	// Archive and prepare new cycle
	archivedTables, err := utils.PrepareNewCycle(oc.DB, archiveDate)
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to archive cycle: %v", err))
		return
	}

	response := models.ArchiveResponse{
		Message:        "Cycle archived successfully and new cycle prepared",
		ArchivedTables: archivedTables,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}
