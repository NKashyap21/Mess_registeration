package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/services"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (a *AuthController) GoogleLoginHandler(c *gin.Context) {
	logger := services.GetLoggerService()
	var payload models.GoogleLoginData

	// Check Content-Type header
	if c.Request.Header.Get("Content-Type") != "application/json" {
		utils.RespondWithError(c, http.StatusBadRequest, "Content-Type must be application/json")
		logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid Content-Type header", c.ClientIP())
		return
	}

	// Decode the request body into the payload struct
	err := utils.ParseJSONRequest(c, &payload)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("Invalid request payload: %s", err.Error()), c.ClientIP())
		return
	}

	// Check if the token is empty after trimming whitespace
	if strings.TrimSpace(payload.Token) == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Token is required")
		logger.LogAuthAction(0, "LOGIN_FAILED", "Empty token provided", c.ClientIP())
		return
	}
	token := payload.Token

	idInfo, err := a.verifyGoogleIDToken(token)
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("Invalid Google token: %s", err.Error()), c.ClientIP())
		return
	}

	email, ok := idInfo["email"].(string)
	if !ok || email == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid token payload")
		logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid email in token payload", c.ClientIP())
		return
	}

	// Find user by encrypted email
	var user models.User
	err = a.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not found")
		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("User not found for email: %s", email), c.ClientIP())
		return
	}

	// Generate JWT token
	var tokenString string
	if tokenString, err = services.GenerateJWT(user.ID, config.GetJWTConfig().SecretKey); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Error creating token")
		logger.LogAuthAction(user.ID, "LOGIN_FAILED", "Error generating JWT token", c.ClientIP())
		return
	}

	// Log successful login
	logger.LogAuthAction(user.ID, "LOGIN_SUCCESS", fmt.Sprintf("User %s logged in successfully", email), c.ClientIP())

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Login successful",
		Data:    gin.H{"token": tokenString, "user": user},
	})
}

// Function to verify Google ID token
func (a *AuthController) verifyGoogleIDToken(token string) (map[string]interface{}, error) {
	req, err := utils.CreateHTTPRequest("https://oauth2.googleapis.com/tokeninfo?id_token="+token, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.SendHTTPRequest(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var idInfo map[string]interface{}
	if err := utils.ParseJSONResponse(resp, &idInfo); err != nil {
		return nil, err
	}

	return idInfo, nil
}
