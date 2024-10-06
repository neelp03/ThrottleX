package test

import (
    "context"
    "github.com/go-redis/redis/v8"
    "testing"
    "time"
    "github.com/neelp03/throttlex/pkg/ratelimiter"
)

func TestFixedWindowLimiter(t *testing.T) {
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Call NewFixedWindowLimiter using the full package name
    limiter := ratelimiter.NewFixedWindowLimiter(redisClient, 5, time.Minute)

    key := "test-key"

    // First 5 requests should be allowed
    for i := 0; i < 5; i++ {
        allowed, err := limiter.Allow(context.Background(), key)
        if err != nil || !allowed {
            t.Errorf("Expected request to be allowed, but got denied on attempt %d", i+1)
        }
    }

    // 6th request should be denied
    allowed, err := limiter.Allow(context.Background(), key)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if allowed {
        t.Error("Expected request to be denied, but it was allowed")
    }
}
