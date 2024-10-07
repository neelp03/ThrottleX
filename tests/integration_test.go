// tests/integration_test.go

package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/neelp03/throttlex/ratelimiter"
	"github.com/neelp03/throttlex/store"
)

func TestFixedWindowLimiter_MemoryStore_Concurrent(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, nil := ratelimiter.NewFixedWindowLimiter(memStore, 100, time.Second*1)
	if nil != nil {
		t.Errorf("Failed to create rate limiter: %v", nil)
	}
	key := "user1"

	var wg sync.WaitGroup
	var allowedCount int
	var mu sync.Mutex

	// Simulate 150 concurrent requests
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if allowed {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowedCount != 100 {
		t.Errorf("Expected 100 allowed requests, got %d", allowedCount)
	}
}

func TestFixedWindowLimiter_RedisStore_Concurrent(t *testing.T) {
	// Set up Redis client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := client.Ping(client.Context()).Err()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	redisStore := store.NewRedisStore(client)
	key := "user1"

	// Clean up the key before test
	client.Del(client.Context(), key)

	limiter, err := ratelimiter.NewFixedWindowLimiter(redisStore, 100, time.Second*1)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", nil)
	}

	var wg sync.WaitGroup
	var allowedCount int
	var mu sync.Mutex

	// Simulate 150 concurrent requests
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if allowed {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowedCount != 100 {
		t.Errorf("Expected 100 allowed requests, got %d", allowedCount)
	}

	// Clean up the key after test
	client.Del(client.Context(), key)
}

func TestSlidingWindowLimiter_MemoryStore_Concurrent(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := ratelimiter.NewSlidingWindowLimiter(memStore, 100, time.Second*1)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user1"

	var wg sync.WaitGroup
	var allowedCount int
	var mu sync.Mutex

	// Simulate 150 concurrent requests
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if allowed {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// The allowed count should not exceed 100
	if allowedCount != 100 {
		t.Errorf("Expected 100 allowed requests, got %d", allowedCount)
	}
}

func TestTokenBucketLimiter_MemoryStore_Concurrent(t *testing.T) {
	memStore := store.NewMemoryStore()
	capacity := 100.0  // Maximum tokens
	refillRate := 50.0 // Tokens per second
	limiter, err := ratelimiter.NewTokenBucketLimiter(memStore, capacity, refillRate)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user1"

	var wg sync.WaitGroup
	var allowedCount int
	var mu sync.Mutex

	// Simulate 150 concurrent requests
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if allowed {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Initially, the allowed count should be up to the capacity
	if allowedCount != 100 {
		t.Errorf("Expected 100 allowed requests, got %d", allowedCount)
	}

	// Wait for tokens to refill
	time.Sleep(time.Second * 2)

	// Simulate more requests after refill
	allowedCount = 0
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if allowed {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Should have allowed approximately 100 tokens (refillRate * 2 seconds)
	if allowedCount < 90 || allowedCount > 110 {
		t.Errorf("Expected around 100 allowed requests after refill, got %d", allowedCount)
	}
}
