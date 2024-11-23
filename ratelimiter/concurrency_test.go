package ratelimiter

import (
	"sync"
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

func TestConcurrencyLimiter_AcquireRelease(t *testing.T) {
	s := store.NewMemoryStore()
	cl, err := NewConcurrencyLimiter(s, 2)
	if err != nil {
		t.Fatalf("Error creating ConcurrencyLimiter: %v", err)
	}

	key := "test_concurrency"

	allowed, err := cl.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Expected Allow() to return true")
	}

	allowed, err = cl.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Expected Allow() to return true")
	}

	allowed, err = cl.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if allowed {
		t.Errorf("Expected Allow() to return false when limit reached")
	}

	err = cl.Release(key)
	if err != nil {
		t.Fatalf("Error on Release(): %v", err)
	}

	allowed, err = cl.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Errorf("Expected Allow() to return true after Release()")
	}
}

func TestConcurrencyLimiter_MaxConcurrent(t *testing.T) {
	s := store.NewMemoryStore()
	cl, err := NewConcurrencyLimiter(s, 5)
	if err != nil {
		t.Fatalf("Error creating ConcurrencyLimiter: %v", err)
	}
	key := "test_concurrency"
	var wg sync.WaitGroup
	var maxConcurrent int64
	var mutex sync.Mutex

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := cl.Allow(key)
			if err != nil {
				t.Errorf("Error on Allow(): %v", err)
				return
			}
			if allowed {
				count, err := s.GetCounter(key)
				if err != nil {
					t.Errorf("Error on GetCounter(): %v", err)
					return
				}
				mutex.Lock()
				if count > maxConcurrent {
					maxConcurrent = count
				}
				mutex.Unlock()
				time.Sleep(100 * time.Millisecond)
				err = cl.Release(key)
				if err != nil {
					t.Errorf("Error on Release(): %v", err)
				}
			}
		}()
	}

	wg.Wait()

	if maxConcurrent > 5 {
		t.Errorf("Expected max concurrent to be <= 5, got %d", maxConcurrent)
	}
}

func TestConcurrencyLimiter_Persistence(t *testing.T) {
	s := store.NewMemoryStore()
	cl, err := NewConcurrencyLimiter(s, 5)
	if err != nil {
		t.Fatalf("Error creating ConcurrencyLimiter: %v", err)
	}
	key := "test_concurrency"

	allowed, err := cl.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Fatalf("Expected Allow() to return true")
	}

	allowed, err = cl.Allow(key)
	if err != nil {
		t.Fatalf("Error on Allow(): %v", err)
	}
	if !allowed {
		t.Fatalf("Expected Allow() to return true")
	}

	// Simulate application restart by creating a new limiter with the same key
	cl2, err := NewConcurrencyLimiter(s, 5)
	if err != nil {
		t.Fatalf("Error creating ConcurrencyLimiter: %v", err)
	}
	count, err := s.GetCounter(key)
	if err != nil {
		t.Fatalf("Error on GetCounter(): %v", err)
	}
	if count != 2 {
		t.Errorf("Expected currentConcurrent to be 2, got %d", count)
	}

	err = cl2.Release(key)
	if err != nil {
		t.Fatalf("Error on Release(): %v", err)
	}

	count, err = s.GetCounter(key)
	if err != nil {
		t.Fatalf("Error on GetCounter(): %v", err)
	}
	if count != 1 {
		t.Errorf("Expected currentConcurrent to be 1 after Release, got %d", count)
	}
}
