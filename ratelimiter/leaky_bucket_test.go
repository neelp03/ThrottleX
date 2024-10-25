package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

func TestLeakyBucketLimiter_Allow(t *testing.T) {
	s := store.NewMemoryStore()
	lb, err := NewLeakyBucketLimiter(s, 5, 1.0)
	if err != nil {
		t.Fatalf("Error creating LeakyBucketLimiter: %v", err)
	}

	key := "test_key"

	// Initially, all requests should be allowed until capacity is reached
	for i := 0; i < 5; i++ {
		allowed, err := lb.Allow(key)
		if err != nil {
			t.Fatalf("Error on Allow(): %v", err)
		}
		if !allowed {
			t.Errorf("Expected Allow() to return true at iteration %d", i)
		}
	}

	// Bucket should now be full
	allowed, err := lb.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if allowed {
		t.Errorf("Expected Allow() to return false when bucket is full")
	}

	// Wait for 1 second to leak tokens
	time.Sleep(1 * time.Second)

	// One token should have leaked
	allowed, err = lb.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Expected Allow() to return true after leaking one token")
	}
}

func TestLeakyBucketLimiter_Persistence(t *testing.T) {
	s := store.NewMemoryStore()
	lb, err := NewLeakyBucketLimiter(s, 5, 1.0)
	if err != nil {
		t.Fatalf("Error creating LeakyBucketLimiter: %v", err)
	}

	key := "test_key"

	// Simulate state
	state := &store.LeakyBucketState{
		Queue:        3,
		LastLeakTime: time.Now().Add(-2 * time.Second),
	}
	err = s.SetLeakyBucket(key, state, time.Hour*24)
	if err != nil {
		t.Fatalf("Error setting state: %v", err)
	}

	// Allow should use the persisted state
	allowed, err := lb.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Expected Allow() to return true after loading state")
	}
}

func TestLeakyBucketLimiter_Concurrency(t *testing.T) {
	s := store.NewMemoryStore()
	lb, err := NewLeakyBucketLimiter(s, 100, 100.0)
	if err != nil {
		t.Fatalf("Error creating LeakyBucketLimiter: %v", err)
	}

	key := "test_key"
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			allowed, err := lb.Allow(key)
			if err != nil {
				t.Errorf("Error on Allow(): %v", err)
			}
			if !allowed {
				t.Errorf("Expected Allow() to return true in concurrent execution")
			}
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify that the queue size is 100
	state, err := s.GetLeakyBucket(key)
	if err != nil {
		t.Fatalf("Error getting state: %v", err)
	}
	if state.Queue != 100 {
		t.Errorf("Expected queue size to be 100, got %v", state.Queue)
	}
}
