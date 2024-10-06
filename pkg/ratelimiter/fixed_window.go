package ratelimiter

import (
    "context"
    "github.com/go-redis/redis/v8"
    "time"
)

// FixedWindowLimiter is a rate limiter that uses the fixed window algorithm.
//
// It tracks requests using Redis and allows up to a specified number of requests
// within a given time window. Once the limit is reached, additional requests
// within the same time window are denied.
//
// Fields:
//   - redisClient: The Redis client used for tracking request counts.
//   - limit: The maximum number of requests allowed per time window.
//   - window: The duration of the time window during which requests are counted.
type FixedWindowLimiter struct {
    redisClient *redis.Client
    limit       int           // max requests allowed
    window      time.Duration // duration of the fixed window
}

// NewFixedWindowLimiter initializes a new FixedWindowLimiter instance.
//
// This function creates a new rate limiter using the fixed window algorithm,
// which counts the number of requests within a specified time window.
//
// Params:
//   - redisClient: *redis.Client - The Redis client for tracking request counts.
//   - limit: int - The maximum number of requests allowed within the window.
//   - window: time.Duration - The time window duration for counting requests.
//
// Returns:
//   *FixedWindowLimiter - A new FixedWindowLimiter instance.
//
// Example usage:
//   limiter := ratelimiter.NewFixedWindowLimiter(redisClient, 100, time.Minute)
func NewFixedWindowLimiter(redisClient *redis.Client, limit int, window time.Duration) *FixedWindowLimiter {
    return &FixedWindowLimiter{
        redisClient: redisClient,
        limit:       limit,
        window:      window,
    }
}

// Allow checks if the request is allowed under the current rate limit.
//
// This method increments the request count for the given key (typically an API key or user identifier)
// and determines if the request is allowed within the current fixed time window. If this is the first
// request in the window, it sets the expiration on the Redis key to match the window duration.
// If the limit is exceeded, the request is denied.
//
// Params:
//   - ctx: context.Context - The context for the Redis operation.
//   - key: string - The unique identifier (e.g., API key) for rate limiting.
//
// Returns:
//   bool - Whether the request is allowed (true) or denied (false).
//   error - Any error that occurred while interacting with Redis.
//
// Example usage:
//   allowed, err := limiter.Allow(ctx, "api-key")
//
// Behavior:
//   - If the request count exceeds the limit, it returns false (denied).
//   - If the request is within the limit, it returns true (allowed).
//   - On the first request in the time window, it sets the Redis key's expiration.
func (limiter *FixedWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
    // Construct the Redis key for rate limiting
    redisKey := "ratelimit:" + key

    // Increment the request count for this key
    count, err := limiter.redisClient.Incr(ctx, redisKey).Result()
    if err != nil {
        return false, err
    }

    // Set expiration if this is the first request in the window
    if count == 1 {
        limiter.redisClient.Expire(ctx, redisKey, limiter.window)
    }

    // Deny the request if it exceeds the rate limit
    if count > int64(limiter.limit) {
        return false, nil // request denied
    }

    // Allow the request if it's within the limit
    return true, nil // request allowed
}
