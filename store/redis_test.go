package store

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func setupTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestRedisStore_Increment(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_increment_key"

	// First increment
	count, err := store.Increment(key, 1, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second increment
	count, err = store.Increment(key, 2, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// Negative increment to check for underflow protection
	count, err = store.Increment(key, -5, time.Minute)
	if err != nil {
		t.Fatalf("Increment with negative value failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 after underflow protection, got %d", count)
	}

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_GetCounter(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_get_counter_key"

	// Counter should initially be zero (non-existent key)
	count, err := store.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	// Increment the counter
	_, err = store.Increment(key, 1, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}

	// Counter should be 1 now
	count, err = store.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_AddTimestamp(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_add_timestamp_key"
	timestamp := time.Now().UnixNano()
	expiration := time.Minute

	// Add a timestamp
	err := store.AddTimestamp(key, timestamp, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Verify the timestamp exists
	count, err := store.CountTimestamps(key, timestamp, timestamp)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_CountTimestamps(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_count_timestamps_key"

	now := time.Now().UnixNano()
	err := store.AddTimestamp(key, now, time.Minute)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}
	err = store.AddTimestamp(key, now+1000, time.Minute)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Count timestamps within the range
	count, err := store.CountTimestamps(key, now, now+1000)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Edge case: No timestamps in range
	count, err = store.CountTimestamps(key, now+2000, now+3000)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}

	client.Del(context.Background(), key)
}

func TestRedisStore_TokenBucket(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_token_bucket_key"
	expiration := time.Minute

	// Set token bucket state
	state := &TokenBucketState{
		Tokens:         10,
		LastUpdateTime: time.Now().UnixNano(),
	}
	err := store.SetTokenBucket(key, state, expiration)
	if err != nil {
		t.Fatalf("SetTokenBucket failed: %v", err)
	}

	// Retrieve and check token bucket state
	retrievedState, err := store.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if retrievedState == nil || retrievedState.Tokens != 10 {
		t.Errorf("Expected tokens 10, got %v", retrievedState)
	}

	// Edge case: Non-existent key
	nonExistentKey := "nonexistent_key"
	nonExistentState, err := store.GetTokenBucket(nonExistentKey)
	if err != nil {
		t.Fatalf("Expected no error for non-existent key, got %v", err)
	}
	if nonExistentState != nil {
		t.Errorf("Expected nil state for non-existent key, got %v", nonExistentState)
	}

	client.Del(context.Background(), key)
}

func TestRedisStore_LeakyBucket(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_leaky_bucket_key"
	expiration := time.Minute

	// Set leaky bucket state
	state := &LeakyBucketState{
		Queue:        5,
		LastLeakTime: time.Now(),
	}
	err := store.SetLeakyBucket(key, state, expiration)
	if err != nil {
		t.Fatalf("SetLeakyBucket failed: %v", err)
	}

	// Retrieve and check leaky bucket state
	retrievedState, err := store.GetLeakyBucket(key)
	if err != nil {
		t.Fatalf("GetLeakyBucket failed: %v", err)
	}
	if retrievedState == nil || retrievedState.Queue != 5 {
		t.Errorf("Expected queue 5, got %v", retrievedState)
	}

	// Edge case: Non-existent key
	nonExistentKey := "nonexistent_key_leaky"
	nonExistentLeakyState, err := store.GetLeakyBucket(nonExistentKey)
	if err != nil {
		t.Fatalf("Expected no error for non-existent key, got %v", err)
	}
	if nonExistentLeakyState != nil {
		t.Errorf("Expected nil state for non-existent key, got %v", nonExistentLeakyState)
	}

	client.Del(context.Background(), key)
}
