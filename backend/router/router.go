package router

import (
	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/handlers/auth"
	"github.com/LambdaIITH/mess_registration/handlers/registration"
	"github.com/LambdaIITH/mess_registration/handlers/staff"
	"github.com/LambdaIITH/mess_registration/handlers/status"
	"github.com/LambdaIITH/mess_registration/handlers/swap"
	"github.com/LambdaIITH/mess_registration/handlers/user"
	"github.com/LambdaIITH/mess_registration/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with all routes and middleware
func SetupRouter() *gin.Engine {
	// Initialize router
	r := gin.New()

	// Add global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// API group with /api prefix
	api := r.Group("/api")

	// Initialize controllers
	healthController := status.InitHealthController()
	authController := auth.InitAuthController()
	userController := user.InitUserController()
	registrationController := registration.InitMessController()
	swapController := swap.InitSwapController()
	staffController := staff.InitStaffController()

	// Health check routes
	api.GET("/health", healthController.CheckHealth)
	api.POST("/login", authController.GoogleLoginHandler)

	students := api.Group("/students")
	students.Use(middleware.TokenRequired(config.GetDB(), &gin.Context{}))
	students.GET("/getUser", userController.GetUserInfoHandler)
	students.POST("/registerMess", registrationController.MessRegistrationHandler)
	students.GET("/getSwaps", swapController.GetAllSwapRequestsHandler)
	students.POST("/createSwap", swapController.CreateSwapRequestHandler)
	students.DELETE("/deleteSwap", swapController.DeleteSwapHandler)
	students.POST("/acceptSwap", swapController.AcceptSwapRequestHandler)

	messStaff := api.Group("/messStaff")
	messStaff.GET("/scanning", staffController.ScanningHandler)



	return r
}
