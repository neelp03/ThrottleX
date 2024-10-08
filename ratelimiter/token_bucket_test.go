// ratelimiter/token_bucket_test.go

package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

// TestTokenBucketLimiter tests the TokenBucketLimiter using the MemoryStore.
func TestTokenBucketLimiter(t *testing.T) {
	// Initialize the MemoryStore
	memStore := store.NewMemoryStore()
	capacity := 5.0   // Maximum of 5 tokens
	refillRate := 1.0 // Refill 1 token per second
	limiter, err := NewTokenBucketLimiter(memStore, capacity, refillRate)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user1"

	// Consume all tokens
	for i := 0; i < int(capacity); i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i+1, err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on request exceeding capacity: %v", err)
	}
	if allowed {
		t.Errorf("Request exceeding capacity should not be allowed")
	}

	// Wait for tokens to refill
	time.Sleep(time.Second * 2)

	// Now we should be able to make 2 more requests
	for i := 0; i < 2; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error on refilled request %d: %v", i+1, err)
		}
		if !allowed {
			t.Errorf("Refilled request %d should be allowed", i+1)
		}
	}

	// Next request should be blocked again
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error after consuming refilled tokens: %v", err)
	}
	if allowed {
		t.Errorf("Request after consuming refilled tokens should not be allowed")
	}
}

// TestTokenBucketLimiterEdgeCases tests edge cases for invalid parameters.
func TestTokenBucketLimiterEdgeCases(t *testing.T) {
	memStore := store.NewMemoryStore()

	// Test with negative capacity
	_, err := NewTokenBucketLimiter(memStore, -1, 1.0)
	if err == nil {
		t.Error("Expected error with negative capacity, but got none")
	}

	// Test with zero refill rate
	_, err = NewTokenBucketLimiter(memStore, 5, 0)
	if err == nil {
		t.Error("Expected error with zero refill rate, but got none")
	}

	// Test with negative refill rate
	_, err = NewTokenBucketLimiter(memStore, 5, -1.0)
	if err == nil {
		t.Error("Expected error with negative refill rate, but got none")
	}
}

// TestTokenBucketLimiterRefill verifies that tokens are refilled at the correct rate.
func TestTokenBucketLimiterRefill(t *testing.T) {
	memStore := store.NewMemoryStore()
	capacity := 3.0
	refillRate := 1.0 // 1 token per second
	limiter, err := NewTokenBucketLimiter(memStore, capacity, refillRate)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user2"

	// Consume all tokens
	for i := 0; i < int(capacity); i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i+1, err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on request exceeding capacity: %v", err)
	}
	if allowed {
		t.Errorf("Request exceeding capacity should not be allowed")
	}

	// Wait for 1 token to refill
	time.Sleep(time.Second * 1)

	// One request should now be allowed
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error after refill: %v", err)
	}
	if !allowed {
		t.Errorf("Request after refill should be allowed")
	}

	// Next request should be blocked
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on request after refill: %v", err)
	}
	if allowed {
		t.Errorf("Request should be blocked after single token is consumed")
	}
}
