// ratelimiter/fixed_window_test.go
package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

func TestFixedWindowLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewFixedWindowLimiter(memStore, 5, time.Second*1)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}
	key := "user1"

	// Simulate 5 allowed requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if allowed {
		t.Errorf("6th request should not be allowed")
	}

	// Wait for the window to expire
	time.Sleep(time.Second * 1)

	// Next request should be allowed after window resets
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("Request after window reset should be allowed")
	}

	// Wait for the window to expire
	time.Sleep(time.Second * 1)

	// Edge case: negative limit
	_, err = NewFixedWindowLimiter(memStore, -1, time.Second)
	if err == nil {
		t.Error("Expected error when limit is set to a negative value")
	}

	// Edge case: zero window duration
	_, err = NewFixedWindowLimiter(memStore, 5, 0)
	if err == nil {
		t.Error("Expected error when window duration is set to zero")
	}
}
