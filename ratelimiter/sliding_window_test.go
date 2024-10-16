package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

// TestSlidingWindowLimiterInvalidKeys checks edge cases with invalid key inputs.
func TestSlidingWindowLimiterInvalidKeys(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 5, time.Second*1)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}

	// Empty key
	allowed, err := limiter.Allow("")
	if err == nil || allowed {
		t.Error("Expected error or disallowed access for empty key")
	}

	// Invalid key format
	invalidKey := "invalid!key@format"
	allowed, err = limiter.Allow(invalidKey)
	if err == nil || allowed {
		t.Error("Expected error or disallowed access for invalid key format")
	}
}

// TestSlidingWindowLimiterHighFrequency tests frequent requests within the same window.
func TestSlidingWindowLimiterHighFrequency(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 5, time.Second*2)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}

	key := "highFrequencyUser"

	// Make 5 requests quickly
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i+1, err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// All subsequent requests should be blocked until window partially resets
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on blocked request: %v", err)
	}
	if allowed {
		t.Error("Request should be blocked as the rate limit has been reached")
	}

	// Wait for half of the window duration to pass
	time.Sleep(time.Second)

	// Next request should still be blocked
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on partially reset window: %v", err)
	}
	if allowed {
		t.Error("Request should be blocked, as the partial window reset is not complete")
	}

	// Wait for the rest of the window to expire
	time.Sleep(time.Second)

	// Next request should be allowed after full window reset
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error after window reset: %v", err)
	}
	if !allowed {
		t.Error("Request after window reset should be allowed")
	}
}

// TestSlidingWindowLimiterVariableRequests simulates requests at different intervals.
func TestSlidingWindowLimiterVariableRequests(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 3, time.Second*2)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}

	key := "variableUser"

	// First request - should be allowed
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 1st request: %v", err)
	}
	if !allowed {
		t.Error("1st request should be allowed")
	}

	// Wait a short time and make two more requests within the limit
	time.Sleep(time.Millisecond * 500)
	for i := 0; i < 2; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i+2, err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+2)
		}
	}

	// 4th request should be blocked
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on blocked request: %v", err)
	}
	if allowed {
		t.Error("4th request should be blocked, limit reached")
	}

	// Wait for full window duration and reattempt
	time.Sleep(time.Second * 2)
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error after full window reset: %v", err)
	}
	if !allowed {
		t.Error("Request after window reset should be allowed")
	}
}
