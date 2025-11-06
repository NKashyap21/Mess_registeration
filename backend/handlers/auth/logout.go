package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (a *AuthController) Logout(c *gin.Context) {
	// Clear the JWT cookie for web clients
	c.SetCookie("jwt", "", -1, "/", strings.Split(os.Getenv("FRONTEND_URL"), ":")[1][2:], false, true)

	// Return JSON response for mobile clients
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Logout successful",
		Data:    nil,
	})
}
