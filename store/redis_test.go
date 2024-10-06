// store/redis_test.go

package store

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func setupRedisStore(t *testing.T) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	err := client.Ping(client.Context()).Err()
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	return NewRedisStore(client)
}

func TestRedisStoreIncrement(t *testing.T) {
	store := setupRedisStore(t)
	key := "test_increment_key"
	expiration := time.Second * 10

	// Ensure key is deleted before test
	store.client.Del(store.ctx, key)

	// First increment
	count, err := store.Increment(key, expiration)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second increment
	count, err = store.Increment(key, expiration)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Cleanup
	store.client.Del(store.ctx, key)
}

func TestRedisStoreAddTimestampAndCount(t *testing.T) {
	store := setupRedisStore(t)
	key := "test_timestamps_key"
	expiration := time.Second * 10

	// Ensure key is deleted before test
	store.client.Del(store.ctx, key)

	now := time.Now().UnixNano()

	// Add timestamps
	err := store.AddTimestamp(key, now, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}
	err = store.AddTimestamp(key, now+1, expiration)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	// Count timestamps
	count, err := store.CountTimestamps(key, now, now+1)
	if err != nil {
		t.Fatalf("CountTimestamps failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Cleanup
	store.client.Del(store.ctx, key)
}

func TestRedisStoreTokenBucket(t *testing.T) {
	store := setupRedisStore(t)
	key := "test_token_bucket_key"
	expiration := time.Hour * 1

	// Ensure key is deleted before test
	store.client.Del(store.ctx, key)

	// Initial state should be nil
	state, err := store.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if state != nil {
		t.Errorf("Expected state to be nil, got %+v", state)
	}

	// Set token bucket state
	initialState := &TokenBucketState{
		Tokens:         5.0,
		LastUpdateTime: time.Now().UnixNano(),
	}
	err = store.SetTokenBucket(key, initialState, expiration)
	if err != nil {
		t.Fatalf("SetTokenBucket failed: %v", err)
	}

	// Retrieve token bucket state
	state, err = store.GetTokenBucket(key)
	if err != nil {
		t.Fatalf("GetTokenBucket failed: %v", err)
	}
	if state == nil {
		t.Fatalf("Expected state to be not nil")
	}
	if state.Tokens != 5.0 {
		t.Errorf("Expected Tokens to be 5.0, got %f", state.Tokens)
	}

	// Cleanup
	store.client.Del(store.ctx, key)
}
