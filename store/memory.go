package store

import (
	"sync"
	"time"
)

// MemoryStore is an in-memory implementation of the Store interface.
type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]*entry
}

type entry struct {
	count       int64
	expiration  time.Time
	timestamps  []int64           // For sliding window
	tokenBucket *TokenBucketState // For token bucket algorithm
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		data: make(map[string]*entry),
	}
	go store.startCleanup()
	return store
}

// startCleanup runs a background goroutine to remove expired entries.
func (m *MemoryStore) startCleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	for range ticker.C {
		m.mu.Lock()
		now := time.Now()
		for key, e := range m.data {
			if now.After(e.expiration) {
				delete(m.data, key)
			}
		}
		m.mu.Unlock()
	}
}

// Increment increments the counter for the given key by 1.
func (m *MemoryStore) Increment(key string, expiration time.Duration) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	e, exists := m.data[key]

	if !exists || now.After(e.expiration) {
		// Initialize the counter for a new or expired key
		m.data[key] = &entry{
			count:      1,
			expiration: now.Add(expiration),
		}
		return 1, nil
	}

	// Increment the counter for existing key
	e.count++
	e.expiration = now.Add(expiration) // Refresh expiration
	return e.count, nil
}

// AddTimestamp adds a timestamp to the list associated with the key.
func (m *MemoryStore) AddTimestamp(key string, timestamp int64, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	e, exists := m.data[key]

	if !exists || now.After(e.expiration) {
		// Initialize a new entry
		m.data[key] = &entry{
			timestamps: []int64{timestamp},
			expiration: now.Add(expiration),
		}
		return nil
	}

	e.timestamps = append(e.timestamps, timestamp)
	e.expiration = now.Add(expiration) // Refresh expiration
	return nil
}

// CountTimestamps counts the number of timestamps within the given range [start, end].
func (m *MemoryStore) CountTimestamps(key string, start int64, end int64) (int64, error) {
	m.mu.RLock()
	e, exists := m.data[key]
	m.mu.RUnlock()

	now := time.Now()
	if !exists || now.After(e.expiration) {
		// Entry does not exist or is expired
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()
		return 0, nil
	}

	// Filter timestamps within the range
	var count int64
	var validTimestamps []int64
	for _, ts := range e.timestamps {
		if ts >= start && ts <= end {
			count++
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// Update timestamps with valid ones
	m.mu.Lock()
	e.timestamps = validTimestamps
	m.mu.Unlock()

	return count, nil
}

// GetTokenBucket retrieves the current state of the token bucket.
func (m *MemoryStore) GetTokenBucket(key string) (*TokenBucketState, error) {
	m.mu.RLock()
	e, exists := m.data[key]
	m.mu.RUnlock()

	now := time.Now()
	if !exists || now.After(e.expiration) {
		// Entry does not exist or is expired
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()
		return nil, nil
	}

	return e.tokenBucket, nil
}

// SetTokenBucket updates the state of the token bucket.
func (m *MemoryStore) SetTokenBucket(key string, state *TokenBucketState, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	e, exists := m.data[key]

	if !exists || now.After(e.expiration) {
		// Initialize a new entry
		m.data[key] = &entry{
			tokenBucket: state,
			expiration:  now.Add(expiration),
		}
		return nil
	}

	e.tokenBucket = state
	e.expiration = now.Add(expiration) // Refresh expiration
	return nil
}
