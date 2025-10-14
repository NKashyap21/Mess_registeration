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

// TODO: state param for csrf
func (a *AuthController) GoogleLoginRedirect(c *gin.Context) {
	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "Redirect Url",
		Data:    gin.H{"redirect": fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?response_type=code&client_id=%s&scope=openid%%20profile%%20email&redirect_uri=%s", os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("BACKEND_URL")+"/api/login-code")},
	})
}

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
	if tokenString, err = services.GenerateJWT(user.ID, email.(string), jwtData["picture"].(string), jwtData["name"].(string), config.GetJWTConfig().SecretKey); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Error creating token")
		logger.LogAuthAction(user.ID, "LOGIN_FAILED", "Error generating JWT token", c.ClientIP())
		return
	}
	//
	// // Log successful login
	logger.LogAuthAction(user.ID, "LOGIN_SUCCESS", fmt.Sprintf("User %s logged in successfully", email), c.ClientIP())

	c.SetCookie("jwt", tokenString, int(jwtData["exp"].(float64)-jwtData["iat"].(float64)), "/", os.Getenv("FRONTEND_URL"), false, true)
	c.Redirect(303, os.Getenv("FRONTEND_URL"))

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

// Function to verify Google ID token
func (a *AuthController) verifyGoogleIDToken(token string) (map[string]interface{}, error) {
	req, err := utils.CreateHTTPRequest("https://oauth2.googleapis.com/token?code="+token, "POST", nil, nil)
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
