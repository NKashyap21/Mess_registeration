package controller

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/internal/services"
	"github.com/LambdaIITH/mess_registration/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwapController struct {
	swapService *services.SwapService
}

func NewSwapController(db *gorm.DB) *SwapController {
	return &SwapController{
		swapService: services.NewSwapService(db),
	}
}

// CreateSwapRequest handles POST /api/swap-request
func (sc *SwapController) CreateSwapRequest(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	id, exists := c.Get("id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := id.(uuid.UUID)

	var req struct {
		Name   string `json:"name" binding:"required"`
		Email  string `json:"email" binding:"required,email"`
		Type   string `json:"type" binding:"required,oneof=friend public"`
		UserID string `json:"user_id"` // This seems redundant as we get it from JWT
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if err := sc.swapService.CreateSwapRequest(userID, req.Name, req.Email, req.Type); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create swap request", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Swap request created successfully", nil)
}

// DeleteSwapRequest handles DELETE /api/swap-request
func (sc *SwapController) DeleteSwapRequest(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	id, exists := c.Get("id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := id.(uuid.UUID)

	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	if err := sc.swapService.DeleteSwapRequest(userID, req.Name, req.Email); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete swap request", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Swap request deleted successfully", nil)
}

// GetSwaps handles GET /api/get-swaps
func (sc *SwapController) GetSwaps(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	id, exists := c.Get("id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	userID := id.(uuid.UUID)

	swaps, err := sc.swapService.GetSwapRequestsByUser(userID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Swaps retrieved successfully", swaps)
}

// AutoMatchSwaps handles POST /api/auto-match (admin only)
func (sc *SwapController) AutoMatchSwaps(c *gin.Context) {
	if err := sc.swapService.AutoMatchSwaps(); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Auto matching completed successfully", nil)
}
