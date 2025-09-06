package main

import (
	"fmt"
	"os"

	"github.com/LambdaIITH/mess_registration/config"
	"github.com/LambdaIITH/mess_registration/internal/router"
)

func init() {
	config.LoadEnvVaariables()
	// Load config.json
	config.LoadConfig()

	// Connect DB
	config.ConnectDB()
}

func main() {
	PORT := ":" + os.Getenv("PORT")

	r := router.SetupRouter()
	if err := r.Run(PORT); err != nil {
		fmt.Println("Error starting Mess Registration backend:", err)
	}
}
