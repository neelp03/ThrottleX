package ratelimiter

import (
    "context"
    "github.com/go-redis/redis/v8"
    "time"
)

type FixedWindowLimiter struct {
    redisClient *redis.Client
    limit       int           // max requests allowed
    window      time.Duration // duration of the fixed window
}

// NewFixedWindowLimiter creates a new FixedWindowLimiter
func NewFixedWindowLimiter(redisClient *redis.Client, limit int, window time.Duration) *FixedWindowLimiter {
    return &FixedWindowLimiter{
        redisClient: redisClient,
        limit:       limit,
        window:      window,
    }
}

// Allow checks if the request is allowed under the rate limit
func (limiter *FixedWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
    redisKey := "ratelimit:" + key
    count, err := limiter.redisClient.Incr(ctx, redisKey).Result()
    if err != nil {
        return false, err
    }

    if count == 1 {
        // Set the expiration when it's the first request in the window
        limiter.redisClient.Expire(ctx, redisKey, limiter.window)
    }

    if count > int64(limiter.limit) {
        return false, nil // request denied
    }

    return true, nil // request allowed
}
