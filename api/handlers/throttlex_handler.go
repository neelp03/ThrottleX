package handlers

import (
    "context"
    "net/http"
    "github.com/neelp03/throttlex/pkg/ratelimiter"
    "github.com/gin-gonic/gin"
    "github.com/neelp03/throttlex/pkg/database"
    "time"
)

// limiter is a rate limiter instance initialized with a fixed window algorithm.
// It allows up to 100 requests per minute, stored in Redis for efficient management.
var limiter = ratelimiter.NewFixedWindowLimiter(database.RedisClient, 100, time.Minute)

// CheckRateLimit handles GET requests to the /api/throttlex/check endpoint.
//
// It performs rate limiting based on the API key provided in the "X-API-KEY" header.
// The function checks whether the incoming request exceeds the allowed rate limits
// using a fixed window rate-limiting strategy. If the limit is exceeded, it returns
// a 429 (Too Many Requests) HTTP response. If the request is within the allowed
// limit, it responds with a 200 (OK) status.
//
// Rate-limiting is enforced using Redis as the underlying storage mechanism, ensuring
// scalability for distributed environments.
//
// Params:
//   c: *gin.Context - The Gin context that holds the HTTP request/response details.
//
// Example usage:
//   curl -H "X-API-KEY: your_api_key" -X GET http://yourapi.com/api/throttlex/check
//
// Responses:
//   200 OK: "Request allowed" - When the request is within the rate limit.
//   429 Too Many Requests: "Rate limit exceeded" - When the rate limit is exceeded.
//   500 Internal Server Error: "Internal server error" - In case of unexpected issues with rate-limiting.
func CheckRateLimit(c *gin.Context) {
    // Retrieve the API key from the request header "X-API-KEY"
    apiKey := c.GetHeader("X-API-KEY")

    // Perform rate-limiting check using the provided API key
    allowed, err := limiter.Allow(context.Background(), apiKey)
    if err != nil {
        // Return a 500 error if an issue occurs during rate-limiting check
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    if !allowed {
        // Return a 429 error if the rate limit has been exceeded
        c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
        return
    }

    // Return a 200 status if the request is allowed
    c.JSON(http.StatusOK, gin.H{"message": "Request allowed"})
}
