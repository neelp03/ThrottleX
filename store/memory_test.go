package store

import (
	"testing"
	"time"
)

func TestMemoryStore_Increment(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_key"
	expiration := time.Second

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

	// Wait for expiration and check count reset
	time.Sleep(expiration)
	count, err = memStore.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after expiration, got %d", count)
	}
}

func TestMemoryStore_AddTimestamp(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_key"
	timestamp := time.Now().Unix()
	expiration := 100 * time.Millisecond

	// Add a timestamp with cleanup
	err := memStore.addTimestampWithCleanup(key, timestamp, expiration, true)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Verify timestamp exists
	count, err := memStore.CountTimestamps(key, timestamp, timestamp)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Wait for expiration and verify cleanup
	time.Sleep(expiration + 10*time.Millisecond)
	count, err = memStore.CountTimestamps(key, timestamp, timestamp)
	if err != nil {
		t.Fatalf("CountTimestamps failed after cleanup: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after expiration, got %d", count)
	}
}

func TestMemoryStore_CountTimestamps(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_key"
	start := time.Now().Unix()
	timestamp1 := start + 1
	timestamp2 := start + 2
	expiration := time.Minute

	// Add timestamps
	_ = memStore.AddTimestamp(key, timestamp1, expiration)
	_ = memStore.AddTimestamp(key, timestamp2, expiration)

	// Verify timestamps count within range
	count, err := memStore.CountTimestamps(key, start, timestamp2)
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
	expiration := time.Second

	// Set token bucket state
	tokenState := &TokenBucketState{Tokens: 10, LastUpdateTime: time.Now().UnixNano()}
	err := memStore.SetTokenBucket(key, tokenState, expiration)
	if err != nil {
		t.Fatalf("SetTokenBucket failed: %v", err)
	}

	// Get token bucket state
	state, err := memStore.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if state.Tokens != tokenState.Tokens {
		t.Errorf("Expected tokens %f, got %f", tokenState.Tokens, state.Tokens)
	}

	// Wait for expiration and verify cleanup
	time.Sleep(expiration + 10*time.Millisecond)
	state, err = memStore.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed after expiration: %v", err)
	}
	if state != nil {
		t.Errorf("Expected nil state after expiration, got %v", state)
	}
}

func TestMemoryStore_LeakyBucket(t *testing.T) {
	memStore := NewMemoryStore()
	key := "test_leaky_bucket"
	expiration := time.Second

	// Set leaky bucket state
	leakyState := &LeakyBucketState{Queue: 5, LastLeakTime: time.Now()}
	err := memStore.SetLeakyBucket(key, leakyState, expiration)
	if err != nil {
		t.Fatalf("SetLeakyBucket failed: %v", err)
	}

	// Get leaky bucket state
	state, err := memStore.GetLeakyBucket(key)
	if err != nil {
		t.Fatalf("GetLeakyBucket failed: %v", err)
	}
	if state.Queue != leakyState.Queue {
		t.Errorf("Expected queue %d, got %d", leakyState.Queue, state.Queue)
	}

	// Wait for expiration and verify cleanup
	time.Sleep(expiration + 10*time.Millisecond)
	state, err = memStore.GetLeakyBucket(key)
	if err != nil {
		t.Fatalf("GetLeakyBucket failed after expiration: %v", err)
	}
	if state != nil {
		t.Errorf("Expected nil state after expiration, got %v", state)
	}
}
