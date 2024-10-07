// ratelimiter/sliding_window_test.go

package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

// TestSlidingWindowLimiter tests the SlidingWindowLimiter using the MemoryStore.
func TestSlidingWindowLimiter(t *testing.T) {
	// Initialize the MemoryStore
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 5, time.Second*1)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user1"

	// Simulate 5 allowed requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i+1, err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 6th request: %v", err)
	}
	if allowed {
		t.Errorf("6th request should not be allowed")
	}

	// Wait for half the window to pass
	time.Sleep(time.Millisecond * 500)

	// 7th request should still be blocked
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 7th request: %v", err)
	}
	if allowed {
		t.Errorf("7th request should not be allowed")
	}

	// Wait for the window to expire
	time.Sleep(time.Millisecond * 600)

	// Next request should be allowed after window resets
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error after window reset: %v", err)
	}
	if !allowed {
		t.Errorf("Request after window reset should be allowed")
	}
}
