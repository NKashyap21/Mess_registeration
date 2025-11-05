package hosteloffice

import (
	"log"
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// ListArchivedTables lists all archived tables
func (oc *OfficeController) ListArchivedTables(c *gin.Context) {
	usersTables, err := utils.ListArchivedTables(oc.DB, "users")
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list archived users tables")
		return
	}

	scansTables, err := utils.ListArchivedTables(oc.DB, "scans")
	if err != nil {
		log.Println(err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to list archived scans tables")
		return
	}

	response := models.ArchivedTablesListResponse{
		Users: usersTables,
		Scans: scansTables,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}
