// store/memory_test.go

package store

import (
	"testing"
	"time"
)

// TestMemoryStoreIncrement tests the Increment method of MemoryStore.
func TestMemoryStoreIncrement(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_increment"
	expiration := time.Second * 1

	// First increment
	count, err := memStore.Increment(key, expiration)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second increment
	count, err = memStore.Increment(key, expiration)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Wait for expiration
	time.Sleep(expiration + time.Millisecond*100)

	// Counter should reset
	count, err = memStore.Increment(key, expiration)
	if err != nil {
		t.Fatalf("Increment after expiration failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count to reset to 1 after expiration, got %d", count)
	}
}

// TestMemoryStoreAddTimestampAndCount tests AddTimestamp and CountTimestamps methods.
func TestMemoryStoreAddTimestampAndCount(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_timestamps"
	expiration := time.Second * 2
	now := time.Now().UnixNano()

	// Add timestamps
	err := memStore.AddTimestamp(key, now, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}
	err = memStore.AddTimestamp(key, now+1, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Count timestamps within range
	count, err := memStore.CountTimestamps(key, now, now+1)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Wait for timestamps to expire
	time.Sleep(expiration + time.Millisecond*100)

	// Count should be zero after expiration
	count, err = memStore.CountTimestamps(key, now, now+expiration.Nanoseconds())
	if err != nil {
		t.Fatalf("CountTimestamps after expiration failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after expiration, got %d", count)
	}
}

// TestMemoryStoreTokenBucket tests GetTokenBucket and SetTokenBucket methods.
func TestMemoryStoreTokenBucket(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_token_bucket"
	expiration := time.Second * 2
	now := time.Now().UnixNano()

	// Initial state should be nil
	state, err := memStore.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if state != nil {
		t.Errorf("Expected initial state to be nil, got %+v", state)
	}

	// Set token bucket state
	initialState := &TokenBucketState{
		Tokens:         5.0,
		LastUpdateTime: now,
	}
	err = memStore.SetTokenBucket(key, initialState, expiration)
	if err != nil {
		t.Fatalf("SetTokenBucket failed: %v", err)
	}

	// Retrieve token bucket state
	state, err = memStore.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if state == nil {
		t.Fatalf("Expected state to be not nil")
	}
	if state.Tokens != 5.0 {
		t.Errorf("Expected Tokens to be 5.0, got %f", state.Tokens)
	}

	// Wait for expiration
	time.Sleep(expiration + time.Millisecond*100)

	// State should be nil after expiration
	state, err = memStore.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket after expiration failed: %v", err)
	}
	if state != nil {
		t.Errorf("Expected state to be nil after expiration, got %+v", state)
	}
}
