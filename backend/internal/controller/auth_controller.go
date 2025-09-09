package controller

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/internal/schema"
	"github.com/LambdaIITH/mess_registration/internal/services"
	"github.com/LambdaIITH/mess_registration/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		userService: services.NewUserService(db),
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	RollNo   string `json:"rollno"`
	UserType string `json:"user_type"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token    string `json:"token"`
	UserType string `json:"user_type"`
	Message  string `json:"message"`
}

// Login handles user authentication and JWT token generation
// POST /api/login
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Sanitize inputs
	req.Email = utils.SanitizeEmail(req.Email)
	req.Name = utils.SanitizeName(req.Name)
	if req.RollNo != "" {
		req.RollNo = utils.SanitizeRollNumber(req.RollNo)
	}

	// Check if user exists
	user, err := ac.userService.GetUserByUserID(req.UserID)
	if err != nil {
		// User doesn't exist, create new user (first-time login)
		if req.RollNo == "" {
			utils.BadRequestResponse(c, "Roll number is required for first-time login")
			return
		}

		// Validate roll number format
		if err := utils.ValidateRollNumber(req.RollNo); err != nil {
			utils.BadRequestResponse(c, err.Error())
			return
		}

		// Set default user type if not provided
		if req.UserType == "" {
			req.UserType = "student"
		}

		// Validate user type
		if err := utils.ValidateUserType(req.UserType); err != nil {
			utils.BadRequestResponse(c, err.Error())
			return
		}

		newUser := &schema.User{
			UserID:   req.UserID,
			Name:     req.Name,
			Email:    req.Email,
			RollNo:   req.RollNo,
			UserType: req.UserType,
			VegType:  "veg", // Default
			Mess:     0,     // Default mess
		}

		if err := ac.userService.CreateUser(newUser); err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create user", err.Error())
			return
		}

		user = newUser
	} else {
		// User exists, update profile information
		updates := map[string]interface{}{
			"name":  req.Name,
			"email": req.Email,
		}

		// Update roll number if provided and different
		if req.RollNo != "" && req.RollNo != user.RollNo {
			if err := utils.ValidateRollNumber(req.RollNo); err != nil {
				utils.BadRequestResponse(c, err.Error())
				return
			}
			updates["roll_no"] = req.RollNo
		}

		// Update user type if provided and user is admin
		if req.UserType != "" && req.UserType != user.UserType {
			if user.UserType == "admin" || req.UserType == "student" {
				updates["user_type"] = req.UserType
			}
		}

		if err := ac.userService.UpdateUser(user.ID, updates); err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update user profile")
			return
		}

		// Refresh user data
		user, _ = ac.userService.GetUserByID(user.ID)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.UserID, user.ID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", LoginResponse{
		Token:    token,
		UserType: user.UserType,
		Message:  "Authentication successful",
	})
}

// RefreshToken handles JWT token refresh
// POST /api/refresh
func (ac *AuthController) RefreshToken(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	id, exists := c.Get("id")
	if !exists {
		utils.UnauthorizedResponse(c, "Invalid token")
		return
	}

	// Verify user still exists
	user, err := ac.userService.GetUserByUserID(userID.(string))
	if err != nil {
		utils.UnauthorizedResponse(c, "User not found")
		return
	}

	// Generate new JWT token
	token, err := utils.GenerateJWT(user.UserID, user.ID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", gin.H{
		"token": token,
	})
}

// Logout handles user logout (client-side token invalidation)
// POST /api/logout
func (ac *AuthController) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side by removing the token
	// For server-side logout, you would need to implement a token blacklist
	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", gin.H{
		"message": "Please remove the token from client storage",
	})
}

// GetProfile returns the current user's profile
// GET /api/profile
func (ac *AuthController) GetProfile(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	user, err := ac.userService.GetUserByUserID(userID.(string))
	if err != nil {
		utils.UnauthorizedResponse(c, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", gin.H{
		"id":        user.ID,
		"user_id":   user.UserID,
		"name":      user.Name,
		"email":     user.Email,
		"rollno":    user.RollNo,
		"user_type": user.UserType,
		"veg_type":  user.VegType,
		"mess":      user.Mess,
	})
}
