package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// TODO: state param for csrf
func (a *AuthController) GoogleLoginRedirect(c *gin.Context) {
	// Handle mobile POST request with ID token
	if c.Request.Method == "POST" {
		a.handleMobileLogin(c)
		return
	}

	// Original GET request for web OAuth redirect
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Redirect Url",
		Data:    map[string]interface{}{"redirect": fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=%s&scope=openid%%20profile%%20email&redirect_uri=%s", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("BACKEND_URL")+"/api/login-code")},
	})
}