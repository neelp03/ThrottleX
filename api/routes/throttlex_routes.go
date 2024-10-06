package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// RegisterRoutes defines the API routes for the ThrottleX service.
//
// This function registers the routes used in the ThrottleX service, including a
// health check endpoint. The health check endpoint ("/health") is used to verify
// that the service is up and running. It responds with a 200 OK status and a JSON
// message indicating the status of the service.
//
// Params:
//   router: *gin.Engine - The Gin engine that handles HTTP requests.
//
// Example usage:
//   RegisterRoutes(router)
//
// Routes:
//   GET /health - Returns a status message confirming that ThrottleX is operational.
//
// Responses:
//   200 OK: "ThrottleX is running" - Health check success message.
func RegisterRoutes(router *gin.Engine) {
    // Define the health check endpoint
    router.GET("/health", func(c *gin.Context) {
        // Respond with a 200 status and JSON message
        c.JSON(http.StatusOK, gin.H{"status": "ThrottleX is running"})
    })
}
