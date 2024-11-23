package ratelimiter

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

// TestNewSlidingWindowLimiterInvalidParameters tests the creation of the limiter with invalid parameters.
func TestNewSlidingWindowLimiterInvalidParameters(t *testing.T) {
	memStore := store.NewMemoryStore()

	// Test with limit <= 0
	_, err := NewSlidingWindowLimiter(memStore, 0, time.Second*1)
	if err == nil {
		t.Error("Expected error when limit <= 0, got nil")
	}

	// Test with window <= 0
	_, err = NewSlidingWindowLimiter(memStore, 5, time.Duration(0))
	if err == nil {
		t.Error("Expected error when window <= 0, got nil")
	}

	// Test with nil store
	_, err = NewSlidingWindowLimiter(nil, 5, time.Second*1)
	if err == nil {
		t.Error("Expected error when store is nil, got nil")
	}
}

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

// MockStore simulates a store that returns errors for testing error handling.
type MockStore struct{}

func (m *MockStore) Increment(key string, delta int64, expiration time.Duration) (int64, error) {
	return 0, errors.New("mock error")
}
func (m *MockStore) GetCounter(key string) (int64, error) {
	return 0, errors.New("mock error")
}
func (m *MockStore) AddTimestamp(key string, timestamp int64, expiration time.Duration) error {
	return errors.New("mock error")
}
func (m *MockStore) CountTimestamps(key string, start int64, end int64) (int64, error) {
	return 0, errors.New("mock error")
}
func (m *MockStore) GetTokenBucket(key string) (*store.TokenBucketState, error) {
	return nil, errors.New("mock error")
}
func (m *MockStore) SetTokenBucket(key string, state *store.TokenBucketState, expiration time.Duration) error {
	return errors.New("mock error")
}
func (m *MockStore) GetLeakyBucket(key string) (*store.LeakyBucketState, error) {
	return nil, errors.New("mock error")
}
func (m *MockStore) SetLeakyBucket(key string, state *store.LeakyBucketState, expiration time.Duration) error {
	return errors.New("mock error")
}

// MockStorePartial simulates a store that returns an error on CountTimestamps.
type MockStorePartial struct{}

func (m *MockStorePartial) AddTimestamp(key string, timestamp int64, expiration time.Duration) error {
	return nil
}
func (m *MockStorePartial) CountTimestamps(key string, start int64, end int64) (int64, error) {
	return 0, errors.New("mock error")
}

// Implement other methods as no-ops
func (m *MockStorePartial) Increment(key string, delta int64, expiration time.Duration) (int64, error) {
	return 0, nil
}
func (m *MockStorePartial) GetCounter(key string) (int64, error) {
	return 0, nil
}
func (m *MockStorePartial) GetTokenBucket(key string) (*store.TokenBucketState, error) {
	return nil, nil
}
func (m *MockStorePartial) SetTokenBucket(key string, state *store.TokenBucketState, expiration time.Duration) error {
	return nil
}
func (m *MockStorePartial) GetLeakyBucket(key string) (*store.LeakyBucketState, error) {
	return nil, nil
}
func (m *MockStorePartial) SetLeakyBucket(key string, state *store.LeakyBucketState, expiration time.Duration) error {
	return nil
}

// TestSlidingWindowLimiterStoreErrors tests the limiter's behavior when the store returns errors.
func TestSlidingWindowLimiterStoreErrors(t *testing.T) {
	mockStore := &MockStore{}
	limiter, err := NewSlidingWindowLimiter(mockStore, 5, time.Second*1)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter with mock store: %v", err)
	}

	key := "testKey"

	// Test Allow() when AddTimestamp returns an error
	allowed, err := limiter.Allow(key)
	if err == nil {
		t.Error("Expected error from AddTimestamp, got nil")
	}
	if allowed {
		t.Error("Expected allowed to be false when store returns error")
	}

	// Use MockStorePartial to simulate error on CountTimestamps
	mockStorePartial := &MockStorePartial{}
	limiter, err = NewSlidingWindowLimiter(mockStorePartial, 5, time.Second*1)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter with partial mock store: %v", err)
	}

	allowed, err = limiter.Allow(key)
	if err == nil {
		t.Error("Expected error from CountTimestamps, got nil")
	}
	if allowed {
		t.Error("Expected allowed to be false when store returns error")
	}
}

// NewSlidingWindowLimiterWithCleanupInterval allows setting a custom cleanup interval for testing.
func NewSlidingWindowLimiterWithCleanupInterval(store store.Store, limit int, window time.Duration, cleanupInterval time.Duration) (*SlidingWindowLimiter, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}
	if window <= 0 {
		return nil, errors.New("window duration must be greater than zero")
	}
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	limiter := &SlidingWindowLimiter{
		store:           store,
		limit:           limit,
		window:          window,
		mutexes:         sync.Map{},
		cleanupInterval: cleanupInterval,
		cleanupStopCh:   make(chan struct{}),
	}
	go limiter.startMutexCleanup()
	return limiter, nil
}

// TestSlidingWindowLimiterMutexCleanup tests that mutexes are cleaned up after inactivity.
func TestSlidingWindowLimiterMutexCleanup(t *testing.T) {
	memStore := store.NewMemoryStore()
	cleanupInterval := time.Millisecond * 100
	limiter, err := NewSlidingWindowLimiterWithCleanupInterval(memStore, 5, time.Second*1, cleanupInterval)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}
	defer limiter.StopCleanup()

	key := "testKey"

	// Make a request to ensure the mutex is created
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on first request: %v", err)
	}
	if !allowed {
		t.Error("First request should be allowed")
	}

	// Verify that the mutex exists
	_, exists := limiter.mutexes.Load(key)
	if !exists {
		t.Error("Mutex should exist after request")
	}

	// Wait longer than cleanup interval to allow cleanup to run
	time.Sleep(cleanupInterval * 3)

	// Check if the mutex has been cleaned up
	_, exists = limiter.mutexes.Load(key)
	if exists {
		t.Error("Mutex should have been cleaned up after inactivity")
	}
}

// TestSlidingWindowLimiterStopCleanup tests that the cleanup goroutine stops when requested.
func TestSlidingWindowLimiterStopCleanup(t *testing.T) {
	memStore := store.NewMemoryStore()
	cleanupInterval := time.Millisecond * 100
	limiter, err := NewSlidingWindowLimiterWithCleanupInterval(memStore, 5, time.Second*1, cleanupInterval)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}

	key := "testKey"

	// Make a request to ensure the mutex is created
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on first request: %v", err)
	}
	if !allowed {
		t.Error("First request should be allowed")
	}

	// Verify that the mutex exists
	_, exists := limiter.mutexes.Load(key)
	if !exists {
		t.Error("Mutex should exist after request")
	}

	// Stop the cleanup
	limiter.StopCleanup()

	// Wait longer than cleanup interval to see if cleanup still runs
	time.Sleep(cleanupInterval * 3)

	// Check if the mutex has been cleaned up
	_, exists = limiter.mutexes.Load(key)
	if !exists {
		t.Error("Mutex should not have been cleaned up after StopCleanup")
	}
}

// TestSlidingWindowLimiterConcurrency tests the limiter under concurrent requests.
func TestSlidingWindowLimiterConcurrency(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 100, time.Second*1)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}

	key := "concurrentUser"
	numRequests := 1000
	var wg sync.WaitGroup
	var allowedCount int64
	var errorCount int64

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, err := limiter.Allow(key)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
			}
			if allowed {
				atomic.AddInt64(&allowedCount, 1)
			}
		}()
	}
	wg.Wait()

	if errorCount != 0 {
		t.Errorf("Expected no errors, got %d errors", errorCount)
	}

	// The allowed count should not exceed the limit
	if allowedCount > int64(limiter.limit) {
		t.Errorf("Allowed count %d exceeds limit %d", allowedCount, limiter.limit)
	}
}

// TestSlidingWindowLimiterTimeBoundaries tests requests at the edge of the time window.
func TestSlidingWindowLimiterTimeBoundaries(t *testing.T) {
	memStore := store.NewMemoryStore()
	limiter, err := NewSlidingWindowLimiter(memStore, 2, time.Millisecond*300)
	if err != nil {
		t.Fatalf("Failed to create SlidingWindowLimiter: %v", err)
	}

	key := "boundaryUser"

	// First request at t=0ms
	allowed, err := limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 1st request: %v", err)
	}
	if !allowed {
		t.Error("1st request should be allowed")
	}

	// Second request at t=100ms
	time.Sleep(time.Millisecond * 100)
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 2nd request: %v", err)
	}
	if !allowed {
		t.Error("2nd request should be allowed")
	}

	// Third request at t=200ms
	time.Sleep(time.Millisecond * 100)
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 3rd request: %v", err)
	}
	if allowed {
		t.Error("3rd request should be blocked, limit reached")
	}

	// Wait until t=400ms
	time.Sleep(time.Millisecond * 100)

	// Fourth request at t=400ms
	allowed, err = limiter.Allow(key)
	if err != nil {
		t.Errorf("Unexpected error on 4th request: %v", err)
	}
	if !allowed {
		t.Error("4th request should be allowed, as the first request has expired from the window")
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
