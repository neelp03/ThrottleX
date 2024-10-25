package store

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisStore_Increment(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	store := NewRedisStore(client)
	key := "test_key"

	// First increment
	count, err := store.Increment(key, 1, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second increment
	count, err = store.Increment(key, 1, time.Minute)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}
}
