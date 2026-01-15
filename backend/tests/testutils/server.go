package testutils

import (
	"log"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/migrations"
	"github.com/LambdaIITH/mess_registration/router"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	gin.SetMode(gin.TestMode)

	config.ConnectDatabase()

	log.Println("Running migrations")
	migrations.MigrateDB()

	return router.SetupRouter()
}
