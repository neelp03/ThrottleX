package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neelp03/throttlex/api/routes"
)

// Setup initializes and configures the Gin engine for the ThrottleX API.
//
// This function sets up the default Gin router and registers the necessary
// routes for the ThrottleX API by calling the `RegisterRoutes` function from
// the `routes` package. The router is then returned, ready to handle HTTP
// requests.
//
// This is typically called during the initialization phase of the application
// to set up routing and middleware configurations.
//
// Returns:
//   *gin.Engine - The initialized Gin router configured with the necessary routes.
//
// Example usage:
//   router := api.Setup()
//   router.Run(":8080") // Starts the server on port 8080
func Setup() *gin.Engine {
	// Initialize a new Gin router with default middleware (logging and recovery)
	router := gin.Default()

	// Register application-specific routes
	routes.RegisterRoutes(router)

	// Return the configured router
	return router
}
