package test

import (
    "context"
    "github.com/go-redis/redis/v8"
    "testing"
    "time"
    "github.com/neelp03/throttlex/pkg/ratelimiter"
)

// TestFixedWindowLimiter validates the behavior of the Fixed Window Rate Limiter.
//
// This test sets up a rate limiter that allows 5 requests per minute and ensures
// that the rate limiting logic works as expected. It checks whether the first 5 requests
// are allowed and ensures that the 6th request is denied due to exceeding the limit.
//
// The test is performed using a Redis client connected to localhost, so Redis must be running
// locally for the test to pass.
//
// Steps:
//   - Create a Redis client connected to "localhost:6379".
//   - Initialize a Fixed Window Rate Limiter allowing 5 requests per minute.
//   - Make 5 consecutive requests, all of which should be allowed.
//   - Make a 6th request, which should be denied.
//
// Fatal Errors:
//   The test will terminate early if there are any unexpected errors.
//
// Example usage:
//   go test -v ./test
func TestFixedWindowLimiter(t *testing.T) {
    // Create a new Redis client connected to localhost
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Redis server address
    })

    // Initialize a Fixed Window Rate Limiter allowing 5 requests per minute
    limiter := ratelimiter.NewFixedWindowLimiter(redisClient, 5, time.Minute)

    // Define a test key for rate limiting
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
