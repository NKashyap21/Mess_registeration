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
	api.GET("/login-code", authController.GoogleLoginHandler)
	api.GET("/getUser", middleware.TokenRequired(config.GetDB(), &gin.Context{}), userController.GetUserInfoHandler)
	// api.POST("/login", authController.GoogleLoginRedirect) // For mobile ID token login
	api.POST("/logout", authController.Logout)

	students := api.Group("/students")
	students.Use(middleware.TokenRequired(config.GetDB(), &gin.Context{}))
	// students.GET("/getUser", userController.GetUserInfoHandler)
	students.GET("/isRegistrationOpen", registrationController.IsRegistrationOpen)
	students.GET("/getMess", registrationController.GetUserMessHandler)
	students.GET("/messStats", registrationController.GetMessStatsHandler)
	students.GET("/messStatsGrouped", registrationController.GetMessStatsGroupedHandler)
	students.GET("/getSwaps", swapController.GetAllSwapRequestsHandler)
	students.GET("/getSwapByID", swapController.GetSwapRequestsByID)
	students.POST("/createSwap", swapController.CreateSwapRequestHandler)
	students.POST("/registerMess/:mess", registrationController.MessRegistrationHandler)
	students.POST("/registerVegMess", registrationController.VegMessRegistrationHandler)
	students.POST("/acceptSwap", swapController.AcceptSwapRequestHandler)
	students.DELETE("/deleteSwap", swapController.DeleteSwapHandler)

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
	hostelOffice.GET("/messStatsGrouped", registrationController.GetMessStatsGroupedHandler)
	hostelOffice.GET("/logs", logsController.GetLogsHandler)
	hostelOffice.GET("/logs/user/:user_id", logsController.GetUserActivityHandler)
	hostelOffice.GET("/logs/system", logsController.GetSystemLogsHandler)
	hostelOffice.GET("/logs/stats", logsController.GetLogStatsHandler)
	hostelOffice.GET("/logs/export", logsController.ExportLogsHandler)
	hostelOffice.GET("/logs/range", logsController.GetLogsByDateRangeHandler)
	hostelOffice.POST("/refreshCapacities", registrationController.RefreshCapacitiesHandler)
	hostelOffice.POST("/toggle/reg", officeController.ToggleNormalRegistration)
	hostelOffice.POST("/toggle/veg", officeController.ToggleVegRegistration)
	hostelOffice.PUT("/students/", officeController.EditStudentById)
	hostelOffice.GET("/status", officeController.GetRegistrationStatus)
	hostelOffice.POST("/apply-new-registration", officeController.ApplyNewRegistration)
	hostelOffice.POST("/add-user", officeController.AddNewUser)

	// CSV upload/download endpoints
	hostelOffice.POST("/students/upload-csv", officeController.UploadStudentsCSV)
	hostelOffice.POST("/students/update-can-register-csv", officeController.UpdateCanRegisterCSV)
	hostelOffice.PUT("/students/bulk-edit", officeController.BulkEditStudents)
	hostelOffice.GET("/students/download-csv", officeController.DownloadStudentsCSV)
	hostelOffice.GET("/registrations/download-csv", officeController.DownloadRegistrationsCSV)
	hostelOffice.GET("/scans/download-csv", officeController.DownloadScansCSV)

	// Archive endpoints
	hostelOffice.POST("/archive/cycle", officeController.ArchiveCycle)
	hostelOffice.GET("/archive/list", officeController.ListArchivedTables)
	hostelOffice.GET("/archive/students/download-csv", officeController.DownloadArchivedStudentsCSV)
	hostelOffice.GET("/archive/scans/download-csv", officeController.DownloadArchivedScansCSV)

	return r
}
