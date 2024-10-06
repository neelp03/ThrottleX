package handlers

import (
    "context"
    "net/http"
    "github.com/neelp03/throttlex/pkg/ratelimiter"
    "github.com/gin-gonic/gin"
    "github.com/neelp03/throttlex/pkg/database"
    "time"
)

var limiter = ratelimiter.NewFixedWindowLimiter(database.RedisClient, 100, time.Minute) // 100 requests per minute

// CheckRateLimit handles GET /api/throttlex/check
func CheckRateLimit(c *gin.Context) {
    apiKey := c.GetHeader("X-API-KEY") // API key from request headers

    allowed, err := limiter.Allow(context.Background(), apiKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    if !allowed {
        c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request allowed"})
}
