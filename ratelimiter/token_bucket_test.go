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
	limiter := NewTokenBucketLimiter(memStore, capacity, refillRate)
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
