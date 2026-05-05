package main

import (
	"log"

	"graduateEmploy/config"
	"graduateEmploy/database"
	"graduateEmploy/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Printf("Warning: %v", err)
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	gin.SetMode(config.AppConfig.ServerMode)

	r := routes.SetupRouter()

	log.Printf("Server starting on port %s...", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
