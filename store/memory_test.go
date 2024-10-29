package store

import (
	"testing"
	"time"
)

func TestMemoryStore_Increment(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_key"
	expiration := time.Minute

	// First increment
	count, err := memStore.Increment(key, 1, expiration)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second increment
	count, err = memStore.Increment(key, 1, expiration)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Test decrement
	count, err = memStore.Increment(key, -1, expiration)
	if err != nil {
		t.Fatalf("Decrement failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1 after decrement, got %d", count)
	}

	// Test reset after expiration
	time.Sleep(expiration)
	count, err = memStore.Increment(key, 1, expiration)
	if err != nil {
		t.Fatalf("Increment after expiration failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1 after expiration reset, got %d", count)
	}
}

func TestMemoryStore_GetCounter(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_key"

	// Counter should be 0 initially
	count, err := memStore.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Increment the counter
	_, err = memStore.Increment(key, 1, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}

	// Counter should be 1
	count, err = memStore.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}
}

func TestMemoryStore_AddTimestamp(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_sliding_window"
	timestamp := time.Now().UnixNano()
	expiration := time.Second * 2

	err := memStore.AddTimestamp(key, timestamp, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Wait for expiration and check if the key has been removed
	time.Sleep(expiration + time.Millisecond*100)
	_, exists := memStore.slidingWindows[key]
	if exists {
		t.Error("Expected key to be deleted after expiration, but it still exists")
	}
}

func TestMemoryStore_CountTimestamps(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_sliding_window"

	// Add timestamps
	now := time.Now().UnixNano()
	err := memStore.AddTimestamp(key, now, time.Minute)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}
	err = memStore.AddTimestamp(key, now+1000, time.Minute)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Count timestamps within range
	count, err := memStore.CountTimestamps(key, now, now+1000)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}

func TestMemoryStore_TokenBucket(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_token_bucket"
	expiration := time.Second * 2

	// Set token bucket state
	state := &TokenBucketState{
		Tokens:         10,
		LastUpdateTime: time.Now().UnixNano(),
	}
	err := memStore.SetTokenBucket(key, state, expiration)
	if err != nil {
		t.Fatalf("SetTokenBucket failed: %v", err)
	}

	// Get token bucket state
	retrievedState, err := memStore.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if retrievedState == nil || retrievedState.Tokens != 10 {
		t.Errorf("Expected tokens 10, got %v", retrievedState)
	}

	// Wait for expiration and check if the key has been removed
	time.Sleep(expiration + time.Millisecond*100)
	retrievedState, _ = memStore.GetTokenBucket(key)
	if retrievedState != nil {
		t.Error("Expected token bucket state to be deleted after expiration, but it still exists")
	}
}

func TestMemoryStore_LeakyBucket(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_leaky_bucket"
	expiration := time.Second * 2

	// Set leaky bucket state
	state := &LeakyBucketState{
		Queue:        5,
		LastLeakTime: time.Now(),
	}
	err := memStore.SetLeakyBucket(key, state, expiration)
	if err != nil {
		t.Fatalf("SetLeakyBucket failed: %v", err)
	}

	// Get leaky bucket state
	retrievedState, err := memStore.GetLeakyBucket(key)
	if err != nil {
		t.Fatalf("GetLeakyBucket failed: %v", err)
	}
	if retrievedState == nil || retrievedState.Queue != 5 {
		t.Errorf("Expected queue 5, got %v", retrievedState)
	}

	// Wait for expiration and check if the key has been removed
	time.Sleep(expiration + time.Millisecond*100)
	retrievedState, _ = memStore.GetLeakyBucket(key)
	if retrievedState != nil {
		t.Error("Expected leaky bucket state to be deleted after expiration, but it still exists")
	}
}
