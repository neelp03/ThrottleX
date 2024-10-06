package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neelp03/throttlex/api/routes"
)

// Setup initializes the Gin router and registers routes
func Setup() *gin.Engine {
	router := gin.Default()
	routes.RegisterRoutes(router)
	return router
}
