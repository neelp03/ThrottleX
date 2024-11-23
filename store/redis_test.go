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

func TestRedisStore_Increment_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_increment_key"

	// Normal increment
	count, err := store.Increment(key, 1, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Simulate Redis error
	client.Close()
	_, err = store.Increment(key, 1, time.Minute)
	if err == nil {
		t.Fatalf("Expected Redis error on Increment, got nil")
	}
	client = setupTestRedisClient() // Reconnect for other tests

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_GetCounter_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_get_counter_key"

	// Non-existent key should return 0 without error
	count, err := store.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0 for non-existent key, got %d", count)
	}

	// Set a value and check retrieval
	_, err = store.Increment(key, 2, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	count, err = store.GetCounter(key)
	if err != nil {
		t.Fatalf("GetCounter failed after increment: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Simulate Redis error
	client.Close()
	_, err = store.GetCounter(key)
	if err == nil {
		t.Fatalf("Expected Redis error on GetCounter, got nil")
	}
	client = setupTestRedisClient() // Reconnect

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_AddTimestamp_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_add_timestamp_key"
	timestamp := time.Now().UnixNano()
	expiration := time.Second

	// Add a timestamp
	err := store.AddTimestamp(key, timestamp, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Verify it exists
	count, err := store.CountTimestamps(key, timestamp, timestamp)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Simulate Redis error on AddTimestamp
	client.Close()
	err = store.AddTimestamp(key, timestamp, expiration)
	if err == nil {
		t.Fatalf("Expected Redis error on AddTimestamp, got nil")
	}
	client = setupTestRedisClient()

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_GetTokenBucket_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_token_bucket_key"
	expiration := time.Second

	// Set a token bucket state
	state := &TokenBucketState{Tokens: 5, LastUpdateTime: time.Now().UnixNano()}
	err := store.SetTokenBucket(key, state, expiration)
	if err != nil {
		t.Fatalf("SetTokenBucket failed: %v", err)
	}

	// Retrieve and validate
	retrieved, err := store.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if retrieved == nil || retrieved.Tokens != 5 {
		t.Errorf("Expected tokens 5, got %v", retrieved)
	}

	// Simulate Redis error
	client.Close()
	_, err = store.GetTokenBucket(key)
	if err == nil {
		t.Fatalf("Expected Redis error on GetTokenBucket, got nil")
	}
	client = setupTestRedisClient()

	// Corrupted data test
	client.HSet(context.Background(), key, "tokens", "not_a_number")
	_, err = store.GetTokenBucket(key)
	if err == nil {
		t.Fatalf("Expected error on corrupted data, got nil")
	}
	client.Del(context.Background(), key)
}

func TestRedisStore_SetTokenBucket_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_token_bucket_key"
	expiration := time.Second

	// Simulate Redis error on SetTokenBucket
	client.Close()
	err := store.SetTokenBucket(key, &TokenBucketState{Tokens: 10, LastUpdateTime: time.Now().UnixNano()}, expiration)
	if err == nil {
		t.Fatalf("Expected Redis error on SetTokenBucket, got nil")
	}
	client = setupTestRedisClient()

	// Cleanup
	client.Del(context.Background(), key)
}

func TestRedisStore_GetLeakyBucket_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_leaky_bucket_key"
	expiration := time.Second

	// Set leaky bucket state
	state := &LeakyBucketState{Queue: 3, LastLeakTime: time.Now()}
	err := store.SetLeakyBucket(key, state, expiration)
	if err != nil {
		t.Fatalf("SetLeakyBucket failed: %v", err)
	}

	// Retrieve and validate
	retrieved, err := store.GetLeakyBucket(key)
	if err != nil {
		t.Fatalf("GetLeakyBucket failed: %v", err)
	}
	if retrieved == nil || retrieved.Queue != 3 {
		t.Errorf("Expected queue 3, got %v", retrieved)
	}

	// Simulate Redis error
	client.Close()
	_, err = store.GetLeakyBucket(key)
	if err == nil {
		t.Fatalf("Expected Redis error on GetLeakyBucket, got nil")
	}
	client = setupTestRedisClient()

	// Corrupted data test
	client.HSet(context.Background(), key, "queue", "not_a_number")
	_, err = store.GetLeakyBucket(key)
	if err == nil {
		t.Fatalf("Expected error on corrupted data, got nil")
	}
	client.Del(context.Background(), key)
}

func TestRedisStore_SetLeakyBucket_ErrorHandling(t *testing.T) {
	client := setupTestRedisClient()
	store := NewRedisStore(client)
	key := "test_leaky_bucket_key"
	expiration := time.Second

	// Simulate Redis error on SetLeakyBucket
	client.Close()
	err := store.SetLeakyBucket(key, &LeakyBucketState{Queue: 10, LastLeakTime: time.Now()}, expiration)
	if err == nil {
		t.Fatalf("Expected Redis error on SetLeakyBucket, got nil")
	}
	client = setupTestRedisClient()

	// Cleanup
	client.Del(context.Background(), key)
}
