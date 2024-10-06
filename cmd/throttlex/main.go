package main

import (
	"log"
	"os"

	"github.com/neelp03/throttlex/api"
	"github.com/neelp03/throttlex/config"
	"github.com/neelp03/throttlex/pkg/database"
	"github.com/neelp03/throttlex/pkg/utils"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize logger
	utils.InitializeLogger()

	// Initialize Redis
	database.InitializeRedis()

	// Start API server
	router := api.Setup()
	port := os.Getenv("PORT")
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
