package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/services"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

// func (a *AuthController) handleMobileLogin(c *gin.Context) {
// 	logger := services.GetLoggerService()

// 	var payload struct {
// 		Token string `json:"token"`
// 	}

// 	if err := utils.ParseJSONRequest(c, &payload); err != nil {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid JSON payload", c.ClientIP())
// 		return
// 	}

// 	if payload.Token == "" {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Token is required")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", "Empty token", c.ClientIP())
// 		return
// 	}

// 	parts := strings.Split(payload.Token, ".")
// 	if len(parts) != 3 {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Invalid token format")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid token format", c.ClientIP())
// 		return
// 	}

// 	decoded, err := base64.RawStdEncoding.DecodeString(parts[1])
// 	if err != nil {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Failed to decode token")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", "Token decode error", c.ClientIP())
// 		return
// 	}

// 	var jwtData map[string]any
// 	if err := json.Unmarshal(decoded, &jwtData); err != nil {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Failed to parse token")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", "Token parse error", c.ClientIP())
// 		return
// 	}

// 	email, ok := jwtData["email"].(string)
// 	if !ok || email == "" {
// 		utils.RespondWithError(c, http.StatusBadRequest, "Invalid email in token")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid email", c.ClientIP())
// 		return
// 	}

// 	var user models.User
// 	if err := a.DB.Where("email = ?", email).First(&user).Error; err != nil {
// 		utils.RespondWithError(c, http.StatusUnauthorized, "User not found")
// 		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("User not found: %s", email), c.ClientIP())
// 		return
// 	}

// 	tokenString, err := services.GenerateJWT(user.ID, user.Type, email, jwtData["name"].(string), jwtData["picture"].(string), config.GetJWTConfig().SecretKey)
// 	if err != nil {
// 		utils.RespondWithError(c, http.StatusInternalServerError, "Error creating token")
// 		logger.LogAuthAction(user.ID, "LOGIN_FAILED", "JWT generation error", c.ClientIP())
// 		return
// 	}

// 	logger.LogAuthAction(user.ID, "LOGIN_SUCCESS", fmt.Sprintf("Mobile login: %s", email), c.ClientIP())

// 	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
// 		Message: "Login successful",
// 		Data:    map[string]interface{}{"token": tokenString, "user": user},
// 	})
// }

func (a *AuthController) GoogleLoginHandler(c *gin.Context) {
	logger := services.GetLoggerService()

	code := c.Request.URL.Query().Get("code")
	if code == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Code needs to be present")
		logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid Code", c.ClientIP())
		return
	}
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", os.Getenv("GOOGLE_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GOOGLE_CLIENT_SECRET"))
	data.Set("redirect_uri", os.Getenv("BACKEND_URL")+"/api/login-code")
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(c, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to verify code: "+err.Error())
		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("Failed to verify authorization code: %s", err.Error()), c.ClientIP())
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := utils.SendHTTPRequest(req, nil)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Code verification failed: "+err.Error())
		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("verification of authorization code failed: %s", err.Error()), c.ClientIP())
		return
	}
	data2, _ := io.ReadAll(res.Body)
	var resData map[string]any
	err = json.Unmarshal(data2, &resData)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to parse google returned data")
		logger.LogAuthAction(0, "LOGIN_FAILED", "Failed to parse json data", c.ClientIP())
		return
	}
	idToken := resData["id_token"]
	decoded, _ := base64.RawStdEncoding.DecodeString(strings.Split(idToken.(string), ".")[1])
	var jwtData map[string]any

	json.Unmarshal(decoded, &jwtData)

	email := jwtData["email"]

	var user models.User
	err = a.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not found")
		logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("User not found for email: %s", email), c.ClientIP())
		return
	}

	// Generate JWT token
	var tokenString string
	if tokenString, err = services.GenerateJWT(user.ID, user.Type, email.(string), jwtData["name"].(string), jwtData["picture"].(string), config.GetJWTConfig().SecretKey); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Error creating token")
		logger.LogAuthAction(user.ID, "LOGIN_FAILED", "Error generating JWT token", c.ClientIP())
		return
	}
	//
	// Log successful login
	logger.LogAuthAction(user.ID, "LOGIN_SUCCESS", fmt.Sprintf("User %s logged in successfully", email), c.ClientIP())

	c.SetCookie("mess_jwt", tokenString, int(jwtData["exp"].(float64)-jwtData["iat"].(float64)), "/", strings.Split(os.Getenv("FRONTEND_URL"), "/")[2], false, true)
	c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONTEND_URL"))

	// utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
	// 	Message: "Login successful",
	// 	Data:    gin.H{"token": tokenString, "user": user},
	// })

	// Check Content-Type header
	// if c.Request.Header.Get("Content-Type") != "application/json" {
	// 	utils.RespondWithError(c, http.StatusBadRequest, "Content-Type must be application/json")
	// 	logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid Content-Type header", c.ClientIP())
	// 	return
	// }

	// // Decode the request body into the payload struct
	// err := utils.ParseJSONRequest(c, &payload)
	// if err != nil {
	// 	utils.RespondWithError(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
	// 	logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("Invalid request payload: %s", err.Error()), c.ClientIP())
	// 	return
	// }

	// Check if the token is empty after trimming whitespace
	// if strings.TrimSpace(payload.Code) == "" {
	// 	utils.RespondWithError(c, http.StatusBadRequest, "Code is required")
	// 	logger.LogAuthAction(0, "LOGIN_FAILED", "Empty token provided", c.ClientIP())
	// 	return
	// }
	//
	// idInfo, err := a.verifyGoogleIDToken(token)
	// if err != nil {
	// 	utils.RespondWithError(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
	// 	logger.LogAuthAction(0, "LOGIN_FAILED", fmt.Sprintf("Invalid Google token: %s", err.Error()), c.ClientIP())
	// 	return
	// }
	//
	// email, ok := idInfo["email"].(string)
	// if !ok || email == "" {
	// 	utils.RespondWithError(c, http.StatusBadRequest, "Invalid token payload")
	// 	logger.LogAuthAction(0, "LOGIN_FAILED", "Invalid email in token payload", c.ClientIP())
	// 	return
	// }
	//
	// // Find user by encrypted email
}
