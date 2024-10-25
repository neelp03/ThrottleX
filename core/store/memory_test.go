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
