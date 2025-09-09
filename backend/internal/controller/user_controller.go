package controller

import (
	"net/http"
	"strconv"

	"github.com/LambdaIITH/mess_registration/internal/schema"
	"github.com/LambdaIITH/mess_registration/internal/services"
	"github.com/LambdaIITH/mess_registration/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: services.NewUserService(db),
	}
}

func (uc *UserController) RegisterUser(c *gin.Context) {
	mess := c.Query("mess")
	if mess == "" {
		utils.BadRequestResponse(c, "mess is required")
		return
	}

	messNum, err := strconv.Atoi(mess)

	if err != nil || messNum < 0 || messNum > 3 {
		utils.BadRequestResponse(c, "invalid mess number")
		return
	}

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	user := &schema.User{
		UserID:   req.UserID,
		Name:     "", // Will be filled from OAuth or additional request
		Email:    "", // Will be filled from OAuth or additional request
		RollNo:   "", // Will be filled from OAuth or additional request
		UserType: "student",
		VegType:  "veg", // Default, can be updated later
		Mess:     messNum,
	}

	if err := uc.userService.CreateUser(user); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Registration unsuccessful", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Registration successful", gin.H{
		"user_id": user.UserID,
		"mess":    user.Mess,
	})

}

// GetUser handles GET /api/user
func (uc *UserController) GetUser(c *gin.Context) {
	// get userID from jwt token (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	user, err := uc.userService.GetUserByUserID(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User found", gin.H{
		"userType": user.UserType,
		"name":     user.Name,
		"rollno":   user.RollNo,
		"veg_type": user.VegType,
		"mess":     user.Mess,
		"email":    user.Email,
	})
}

// GetUserByRollNo handles GET /api/user?rollno=xxx
func (uc *UserController) GetUserByRollNo(c *gin.Context) {
	rollno := c.Query("rollno")
	if rollno == "" {
		utils.BadRequestResponse(c, "rollno is required")
		return
	}
	user, err := uc.userService.GetUserByRollNo(rollno)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", gin.H{
		"userType": user.UserType,
		"name":     user.Name,
		"rollno":   user.RollNo,
		"veg_type": user.VegType,
		"mess":     user.Mess,
		"email":    user.Email,
		"id":       user.ID,
		"user_id":  user.UserID,
	})
}

// UpdateUser handles POST /api/update
func (uc *UserController) UpdateUser(c *gin.Context) {
	rollno := c.Query("rollno")
	if rollno == "" {
		utils.BadRequestResponse(c, "rollno is required")
		return
	}

	// Get existing user
	user, err := uc.userService.GetUserByRollNo(rollno)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "User not found", err.Error())
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Update user
	if err := uc.userService.UpdateUser(user.ID, updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Update unsuccessful", err.Error())
		return
	}

	// Get updated user
	updatedUser, _ := uc.userService.GetUserByID(user.ID)

	utils.SuccessResponse(c, http.StatusOK, "success", gin.H{
		"userType": updatedUser.UserType,
		"name":     updatedUser.Name,
		"rollno":   updatedUser.RollNo,
		"veg_type": updatedUser.VegType,
		"mess":     updatedUser.Mess,
		"email":    updatedUser.Email,
	})
}

// GetAllUsers handles GET /api/users (for hostel office)
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", gin.H{
		"users": users,
		"count": len(users),
	})
}

// ScanUser handles GET /api/scanning (for mess staff)
func (uc *UserController) ScanUser(c *gin.Context) {
	rollno := c.Query("rollno")
	if rollno == "" {
		utils.BadRequestResponse(c, "rollno parameter is required")
		return
	}

	user, err := uc.userService.GetUserByRollNo(rollno)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", gin.H{
		"userType": user.UserType,
		"name":     user.Name,
		"rollno":   user.RollNo,
		"veg_type": user.VegType,
		"mess":     user.Mess,
	})
}

// ResetUsers handles DELETE /api/reset (for admin)
func (uc *UserController) ResetUsers(c *gin.Context) {
	if err := uc.userService.ResetAllUsers(); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All user data reset successfully", nil)
}
