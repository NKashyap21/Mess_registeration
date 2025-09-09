package router

import (
	"github.com/LambdaIITH/mess_registration/internal/controller"
	"github.com/LambdaIITH/mess_registration/internal/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Initialize controllers
	userController := controller.NewUserController(db)
	swapController := controller.NewSwapController(db)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Mess Registration API is running",
		})
	})

	// API routes
	api := r.Group("/api")

	// Public routes (no authentication required)
	api.POST("/register", userController.RegisterUser)
	api.POST("/login", authController.Login)

	// Protected routes (require JWT authentication)
	protected := api.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		// Auth routes
		protected.POST("/refresh", authController.RefreshToken)
		protected.POST("/logout", authController.Logout)
		protected.GET("/profile", authController.GetProfile)

		// User routes
		protected.GET("/user", userController.GetUser)

		// Swap routes
		protected.POST("/swap-request", swapController.CreateSwapRequest)
		protected.DELETE("/swap-request", swapController.DeleteSwapRequest)
		protected.GET("/get-swaps", swapController.GetSwaps)
	}

	// Routes with body-based JWT authentication (alternative method)
	bodyAuth := api.Group("/")
	bodyAuth.Use(middlewares.AuthMiddlewareWithBody())
	{
		bodyAuth.GET("/user-alt", userController.GetUser) // Alternative endpoint
	}

	// Mess staff routes (require API key authentication)
	messStaff := api.Group("/")
	messStaff.Use(middlewares.MessAPIKeyMiddleware())
	{
		messStaff.GET("/scanning", userController.ScanUser)
	}

	// Hostel office routes (require JWT + admin check)
	hostelOffice := api.Group("/admin")
	hostelOffice.Use(middlewares.AuthMiddleware(), middlewares.AdminMiddleware(db))
	{
		hostelOffice.GET("/user", userController.GetUserByRollNo)
		hostelOffice.POST("/update", userController.UpdateUser)
		hostelOffice.GET("/users", userController.GetAllUsers)
		hostelOffice.DELETE("/reset", userController.ResetUsers)
		hostelOffice.POST("/auto-match", swapController.AutoMatchSwaps)
	}

	// Development/Testing routes
	if gin.Mode() == gin.DebugMode {
		dev := api.Group("/dev")
		{
			dev.GET("/users", userController.GetAllUsers)
			dev.DELETE("/reset", userController.ResetUsers)
		}
	}

	return r
}
