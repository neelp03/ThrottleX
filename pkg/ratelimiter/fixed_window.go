package ratelimiter

import (
    "context"
    "time"
    "os"
    "strconv"
    "github.com/go-redis/redis/v8"
    "github.com/neelp03/throttlex/pkg/utils"  // For logging
    "github.com/neelp03/throttlex/pkg/metrics" // For Prometheus metrics
)

// FixedWindowLimiter is a rate limiter that implements the fixed window algorithm.
//
// The rate limiter uses Redis to track the number of requests per key within a defined
// time window. If the number of requests exceeds the specified limit within the window,
// subsequent requests are denied until the window resets.
//
// Fields:
//   - redisClient: Redis client for managing request counts.
//   - limit: Maximum allowed requests within the window.
//   - window: Time duration of the rate-limiting window.
type FixedWindowLimiter struct {
    redisClient *redis.Client
    limit       int
    window      time.Duration
}

// NewFixedWindowLimiter initializes a new FixedWindowLimiter.
//
// The rate limit and window duration can be configured via environment variables. If the
// environment variables are not set, default values are used (100 requests per minute).
//
// Environment Variables:
//   - LIMIT: Maximum number of requests allowed (default: 100).
//   - WINDOW: Time window duration in seconds (default: 60 seconds).
//
// Params:
//   - redisClient: *redis.Client - Redis client for managing request counts.
//
// Returns:
//   *FixedWindowLimiter - A new instance of FixedWindowLimiter.
//
// Example usage:
//   limiter := ratelimiter.NewFixedWindowLimiter(redisClient)
func NewFixedWindowLimiter(redisClient *redis.Client) *FixedWindowLimiter {
    // Fetch the rate limit and window size from environment variables, with default values
    limit, err := strconv.Atoi(os.Getenv("LIMIT"))
    if err != nil {
        limit = 100 // Default: 100 requests
    }

    windowSeconds, err := strconv.Atoi(os.Getenv("WINDOW"))
    if err != nil {
        windowSeconds = 60 // Default: 60 seconds
    }
    window := time.Duration(windowSeconds) * time.Second

    return &FixedWindowLimiter{
        redisClient: redisClient,
        limit:       limit,
        window:      window,
    }
}

// Allow checks if the request is allowed under the current rate limit.
//
// This method increments the request count for the provided key and determines
// if the request is allowed within the current fixed time window. If the request
// count exceeds the limit, the request is denied.
//
// It also logs rate limit breaches, sets Redis key expiration, and increments
// Prometheus metrics for monitoring.
//
// Params:
//   - ctx: context.Context - Context for managing request deadlines and cancellation.
//   - key: string - The unique identifier (e.g., API key or user ID) to be rate-limited.
//
// Returns:
//   bool - True if the request is allowed, false if denied due to rate limit breach.
//   error - Any error encountered while interacting with Redis.
//
// Example usage:
//   allowed, err := limiter.Allow(ctx, "api-key")
//
// Behavior:
//   - If the request count exceeds the limit, it returns false (denied).
//   - If the request is within the limit, it returns true (allowed).
func (limiter *FixedWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
    redisKey := "ratelimit:" + key

    // Increment the request count in Redis
    count, err := limiter.redisClient.Incr(ctx, redisKey).Result()
    if err != nil {
        utils.LogError("Failed to increment Redis key", err)
        return false, err
    }

    // Set expiration for the Redis key if this is the first request in the window
    if count == 1 {
        err := limiter.redisClient.Expire(ctx, redisKey, limiter.window).Err()
        if err != nil {
            utils.LogError("Failed to set expiration for Redis key", err)
            return false, err
        }
    }

    // Log the current request count for monitoring purposes
    utils.LogInfo("Request count for key " + key + ": " + strconv.FormatInt(count, 10))

    // Increment the Prometheus total requests metric
    metrics.TotalRequests.Inc()

    // Deny the request if the count exceeds the rate limit
    if count > int64(limiter.limit) {
        utils.LogInfo("Rate limit exceeded for key: " + key)
        metrics.DeniedRequests.Inc() // Increment the denied requests metric
        return false, nil
    }

    // Request allowed
    return true, nil
}
