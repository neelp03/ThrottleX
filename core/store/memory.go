package store

import (
	"sync"
	"time"
)

// MemoryStore is an in-memory implementation of the Store interface.
type MemoryStore struct {
	mu             sync.Mutex
	counters       map[string]*memoryCounter
	slidingWindows map[string][]int64
	tokenBuckets   map[string]*TokenBucketState
	leakyBuckets   map[string]*LeakyBucketState
}

type memoryCounter struct {
	count      int64
	expiration time.Time
}

// NewMemoryStore initializes a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		counters:       make(map[string]*memoryCounter),
		slidingWindows: make(map[string][]int64),
		tokenBuckets:   make(map[string]*TokenBucketState),
		leakyBuckets:   make(map[string]*LeakyBucketState),
	}
}

// Increment increments the counter by delta and sets expiration.
func (s *MemoryStore) Increment(key string, delta int64, expiration time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	counter, exists := s.counters[key]
	if !exists || time.Now().After(counter.expiration) {
		counter = &memoryCounter{
			count:      delta,
			expiration: time.Now().Add(expiration),
		}
		if counter.count < 0 {
			counter.count = 0
		}
		s.counters[key] = counter
	} else {
		counter.count += delta
		if counter.count < 0 {
			counter.count = 0
		}
	}
	return counter.count, nil
}

// GetCounter retrieves the current value of the counter.
func (s *MemoryStore) GetCounter(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	counter, exists := s.counters[key]
	if !exists || time.Now().After(counter.expiration) {
		return 0, nil
	}
	return counter.count, nil
}

// AddTimestamp adds a timestamp to the list associated with the key.
func (s *MemoryStore) AddTimestamp(key string, timestamp int64, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.slidingWindows[key] = append(s.slidingWindows[key], timestamp)

	// Set up a timer to delete the key after expiration
	go func() {
		time.Sleep(expiration)
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.slidingWindows, key)
	}()
	return nil
}

// CountTimestamps counts the number of timestamps within the given range [start, end].
func (s *MemoryStore) CountTimestamps(key string, start int64, end int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamps, exists := s.slidingWindows[key]
	if !exists {
		return 0, nil
	}

	var count int64
	for _, ts := range timestamps {
		if ts >= start && ts <= end {
			count++
		}
	}
	return count, nil
}

// GetTokenBucket retrieves the current state of the token bucket.
func (s *MemoryStore) GetTokenBucket(key string) (*TokenBucketState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.tokenBuckets[key]
	if !exists {
		return nil, nil
	}
	return state, nil
}

// SetTokenBucket updates the state of the token bucket.
func (s *MemoryStore) SetTokenBucket(key string, state *TokenBucketState, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokenBuckets[key] = state

	// Set up a timer to delete the key after expiration
	go func() {
		time.Sleep(expiration)
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.tokenBuckets, key)
	}()
	return nil
}

// GetLeakyBucket retrieves the current state of the leaky bucket.
func (s *MemoryStore) GetLeakyBucket(key string) (*LeakyBucketState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, exists := s.leakyBuckets[key]
	if !exists {
		return nil, nil
	}
	return state, nil
}

// SetLeakyBucket updates the state of the leaky bucket.
func (s *MemoryStore) SetLeakyBucket(key string, state *LeakyBucketState, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.leakyBuckets[key] = state

	// Set up a timer to delete the key after expiration
	go func() {
		time.Sleep(expiration)
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.leakyBuckets, key)
	}()
	return nil
}
