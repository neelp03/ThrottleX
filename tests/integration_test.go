package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/neelp03/throttlex/ratelimiter"
	"github.com/neelp03/throttlex/store"
)

func TestIntegration_FixedWindowLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := ratelimiter.NewFixedWindowLimiter(memStore, 5, time.Second)
	if err != nil {
		t.Fatalf("Error creating FixedWindowLimiter: %v", err)
	}

	key := "user123"

	// Allow 5 requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow(): %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if allowed {
		t.Errorf("6th request should be blocked")
	}

	// Wait for the window to reset
	time.Sleep(time.Second)

	// Next request should be allowed
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Request after window reset should be allowed")
	}
}

func TestIntegration_SlidingWindowLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := ratelimiter.NewSlidingWindowLimiter(memStore, 5, time.Second)
	if err != nil {
		t.Fatalf("Error creating SlidingWindowLimiter: %v", err)
	}

	key := "user123"

	// Allow 5 requests
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow(): %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 6th request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if allowed {
		t.Errorf("6th request should be blocked")
	}

	// Wait for the window to slide past the first request
	time.Sleep(1 * time.Second)

	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Request after window slide should be allowed")
	}
}

func TestIntegration_TokenBucketLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := ratelimiter.NewTokenBucketLimiter(memStore, 5, 1) // Capacity 5 tokens, refill rate 1 token/sec
	if err != nil {
		t.Fatalf("Error creating TokenBucketLimiter: %v", err)
	}

	key := "user123"

	// Consume 5 tokens
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow(): %v", err)
		}
		if !allowed {
			t.Errorf("Token %d should be consumed", i+1)
		}
	}

	// Next request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if allowed {
		t.Errorf("Request should be blocked due to empty token bucket")
	}

	// Wait for 1 second to refill one token
	time.Sleep(time.Second)

	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Request should be allowed after token refill")
	}
}

func TestIntegration_LeakyBucketLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := ratelimiter.NewLeakyBucketLimiter(memStore, 5, 1) // Capacity 5, leak rate 1/sec
	if err != nil {
		t.Fatalf("Error creating LeakyBucketLimiter: %v", err)
	}

	key := "user123"

	// Fill the bucket
	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow(): %v", err)
		}
		if !allowed {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// Next request should be blocked
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if allowed {
		t.Errorf("Request should be blocked due to full bucket")
	}

	// Wait for 1 second to leak one token
	time.Sleep(time.Second)

	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Request should be allowed after token leaked")
	}
}

func TestIntegration_ConcurrencyLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := ratelimiter.NewConcurrencyLimiter(memStore, 2)
	if err != nil {
		t.Fatalf("Error creating ConcurrencyLimiter: %v", err)
	}

	key := "user123"
	var wg sync.WaitGroup
	var mutex sync.Mutex
	activeRequests := 0
	maxActiveRequests := 0

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				t.Errorf("Error on Allow(): %v", err)
				return
			}
			if !allowed {
				t.Logf("Request denied due to concurrency limit")
				return
			}
			// Simulate work
			mutex.Lock()
			activeRequests++
			if activeRequests > maxActiveRequests {
				maxActiveRequests = activeRequests
			}
			mutex.Unlock()

			time.Sleep(100 * time.Millisecond)

			mutex.Lock()
			activeRequests--
			mutex.Unlock()

			err = limiter.Release(key)
			if err != nil {
				t.Errorf("Error on Release(): %v", err)
			}
		}()
	}

	wg.Wait()

	if maxActiveRequests > 2 {
		t.Errorf("Concurrency limit exceeded: max active requests %d", maxActiveRequests)
	}
}

func TestIntegration_MultipleLimiters(t *testing.T) {
	memStore := store.NewMemoryStore()

	// Create multiple limiters
	fixedLimiter, err := ratelimiter.NewFixedWindowLimiter(memStore, 10, time.Second)
	if err != nil {
		t.Fatalf("Error creating FixedWindowLimiter: %v", err)
	}

	tokenLimiter, err := ratelimiter.NewTokenBucketLimiter(memStore, 5, 1)
	if err != nil {
		t.Fatalf("Error creating TokenBucketLimiter: %v", err)
	}

	key := "user123"

	// Simulate requests passing through multiple limiters
	for i := 0; i < 15; i++ {
		allowedFixed, err := fixedLimiter.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow() from FixedWindowLimiter: %v", err)
		}
		allowedToken, err := tokenLimiter.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow() from TokenBucketLimiter: %v", err)
		}

		if allowedFixed && allowedToken {
			t.Logf("Request %d allowed", i+1)
		} else {
			t.Logf("Request %d blocked", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
