package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

// TestSlidingWindowLimiter tests the SlidingWindowLimiter with various scenarios.
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

// TestSlidingWindowLimiterEdgeCases checks edge cases for invalid parameters.
func TestSlidingWindowLimiterEdgeCases(t *testing.T) {
	memStore := store.NewMemoryStore()

	// Test with negative limit
	_, err := NewSlidingWindowLimiter(memStore, -5, time.Second*1)
	if err == nil {
		t.Error("Expected error with negative limit, but got none")
	}

	// Test with zero window duration
	_, err = NewSlidingWindowLimiter(memStore, 5, 0)
	if err == nil {
		t.Error("Expected error with zero window duration, but got none")
	}
}

// TestSlidingWindowLimiterMultipleClients simulates rate limiting for multiple clients.
func TestSlidingWindowLimiterMultipleClients(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 3, time.Second*1)
	if err != nil {
		t.Errorf("Failed to create rate limiter: %v", err)
	}

	// Simulate requests for multiple clients
	keys := []string{"client1", "client2", "client3"}
	for _, key := range keys {
		for i := 0; i < 3; i++ {
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Unexpected error for key %s on request %d: %v", key, i+1, err)
			}
			if !allowed {
				t.Errorf("Request %d for key %s should be allowed", i+1, key)
			}
		}

		// 4th request should be blocked for each key
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error for key %s on 4th request: %v", key, err)
		}
		if allowed {
			t.Errorf("4th request for key %s should not be allowed", key)
		}
	}
}
