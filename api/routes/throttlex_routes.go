package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// RegisterRoutes defines the routes for ThrottleX
func RegisterRoutes(router *gin.Engine) {
    router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ThrottleX is running"})
    })
}
