package test

import (
    "context"
    "testing"
    "time"
    "github.com/go-redis/redismock/v8"  // Mocking Redis interactions
    "github.com/neelp03/throttlex/pkg/ratelimiter"
)

// TestFixedWindowLimiter validates the behavior of the Fixed Window Rate Limiter.
//
// This test uses a Redis mock to simulate interactions with Redis and ensures that
// the rate limiting logic works as expected. The Fixed Window Rate Limiter allows
// up to 5 requests per minute. The test checks whether requests are allowed or denied
// based on the rate limit.
//
// Test cases covered:
//   - First 5 requests are allowed.
//   - The 6th request is denied due to exceeding the limit.
//   - All Redis interactions are mocked using the redismock package.
//
// Redis Mocking:
//   - The INCR command is mocked to return the current request count.
//   - The EXPIRE command is mocked to simulate setting the expiration time on the rate limit key.
//
// Example usage:
//   go test -v ./pkg/ratelimiter/test
func TestFixedWindowLimiter(t *testing.T) {
    // Set up a mock Redis client for simulating Redis interactions.
    redisClient, mock := redismock.NewClientMock()

    // Initialize a Fixed Window Rate Limiter with a limit of 5 requests per minute.
    limiter := ratelimiter.NewFixedWindowLimiter(redisClient, 5, time.Minute)

    key := "test-key"

    // Mock Redis behavior for the first request: INCR should return 1, EXPIRE should set a 1-minute window.
    mock.ExpectIncr("ratelimit:" + key).SetVal(1)
    mock.ExpectExpire("ratelimit:" + key, time.Minute).SetVal(true)

    // First request should be allowed.
    allowed, err := limiter.Allow(context.Background(), key)
    if err != nil || !allowed {
        t.Errorf("Expected first request to be allowed, but got denied or an error: %v", err)
    }

    // Mock Redis behavior for subsequent requests (second to fifth).
    for i := 2; i <= 5; i++ {
        mock.ExpectIncr("ratelimit:" + key).SetVal(int64(i))
    }

    // Next 4 requests (2 to 5) should be allowed.
    for i := 0; i < 4; i++ {
        allowed, err := limiter.Allow(context.Background(), key)
        if err != nil || !allowed {
            t.Errorf("Expected request to be allowed, but got denied on attempt %d", i+2)
        }
    }

    // Mock Redis behavior for the 6th request: INCR should return 6 (exceeding the limit).
    mock.ExpectIncr("ratelimit:" + key).SetVal(6)

    // 6th request should be denied due to exceeding the limit.
    allowed, err = limiter.Allow(context.Background(), key)
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    if allowed {
        t.Error("Expected 6th request to be denied, but it was allowed")
    }

    // Ensure all Redis mock expectations were met.
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("Redis expectations were not met: %v", err)
    }
}
