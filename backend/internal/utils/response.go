package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Message string      `json:"message"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message, error string) {
	c.JSON(statusCode, APIResponse{
		Message: message,
		Error:   error,
	})
}

func BadRequestResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusBadRequest, "Bad request", error)
}

func UnauthorizedResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", error)
}

func InternalServerErrorResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusInternalServerError, "Internal server error", error)
}

func NotFoundResponse(c *gin.Context, error string) {
	ErrorResponse(c, http.StatusNotFound, "Not found", error)
}
