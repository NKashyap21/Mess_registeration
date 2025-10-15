package router

import (
	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/handlers/auth"
	hosteloffice "github.com/LambdaIITH/mess_registration/handlers/hostelOffice"
	"github.com/LambdaIITH/mess_registration/handlers/logs"
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
	r.Use(middleware.Logger())         // Console logging with panic recovery
	r.Use(middleware.DatabaseLogger()) // Database logging
	r.Use(middleware.Recovery())       // Additional recovery middleware

	// API group with /api prefix
	api := r.Group("/api")

	// Initialize controllers
	healthController := status.InitHealthController()
	authController := auth.InitAuthController()
	userController := user.InitUserController()
	registrationController := registration.InitMessController()
	swapController := swap.InitSwapController()
	staffController := staff.InitStaffController()
	officeController := hosteloffice.InitOfficeController()
	logsController := logs.InitLogsController()

	// Health check routes
	api.GET("/health", healthController.CheckHealth)
	api.GET("/login", authController.GoogleLoginRedirect)
	api.POST("/logout", authController.Logout)
	api.GET("/login-code", authController.GoogleLoginHandler)
	api.GET("/getUser", middleware.TokenRequired(config.GetDB(), &gin.Context{}), userController.GetUserInfoHandler)

	students := api.Group("/students")
	students.Use(middleware.TokenRequired(config.GetDB(), &gin.Context{}))
	// students.GET("/getUser", userController.GetUserInfoHandler)
	students.GET("/isRegistrationOpen", registrationController.IsRegistrationOpen)
	students.POST("/registerMess/:mess", registrationController.MessRegistrationHandler)
	students.POST("/registerVegMess", registrationController.VegMessRegistrationHandler)
	students.GET("/getMess", registrationController.GetUserMessHandler)
	students.GET("/messStats", registrationController.GetMessStatsHandler)
	students.GET("/messStatsGrouped", registrationController.GetMessStatsGroupedHandler)
	students.GET("/getSwaps", swapController.GetAllSwapRequestsHandler)
	students.GET("/getSwapByID", swapController.GetSwapRequestsByID)
	students.POST("/createSwap", swapController.CreateSwapRequestHandler)
	students.DELETE("/deleteSwap", swapController.DeleteSwapHandler)
	students.POST("/acceptSwap", swapController.AcceptSwapRequestHandler)

	messStaff := api.Group("/messStaff")
	messStaff.Use(middleware.TokenRequired(config.GetDB(), &gin.Context{}))
	messStaff.Use(middleware.MessStaffMiddleware(config.GetDB()))

	messStaff.GET("/info", staffController.GetStaffInfo)
	messStaff.GET("/scanning", staffController.ScanningHandler)

	hostelOffice := api.Group("/office")
	hostelOffice.Use(middleware.TokenRequired(config.GetDB(), &gin.Context{}))
	hostelOffice.Use(middleware.HostelOfficeMiddleWare(config.GetDB()))

	hostelOffice.GET("/students", officeController.GetStudents)
	hostelOffice.GET("/students/:roll_no", officeController.GetStudentsByID)
	hostelOffice.PUT("/students/", officeController.EditStudentById)
	hostelOffice.POST("/refreshCapacities", registrationController.RefreshCapacitiesHandler)
	hostelOffice.GET("/messStatsGrouped", registrationController.GetMessStatsGroupedHandler)

	hostelOffice.GET("/logs", logsController.GetLogsHandler)
	hostelOffice.GET("/logs/user/:user_id", logsController.GetUserActivityHandler)
	hostelOffice.GET("/logs/system", logsController.GetSystemLogsHandler)
	hostelOffice.GET("/logs/stats", logsController.GetLogStatsHandler)
	hostelOffice.GET("/logs/export", logsController.ExportLogsHandler)
	hostelOffice.GET("/logs/range", logsController.GetLogsByDateRangeHandler)

	return r
}
