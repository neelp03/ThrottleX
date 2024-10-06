package ratelimiter

import (
	"sync"
	"time"

	"github.com/neelp03/throttlex/store"
)

// TokenBucketLimiter implements the token bucket rate-limiting algorithm.
// It allows a certain burst of requests and refills tokens at a steady rate.
type TokenBucketLimiter struct {
	store           store.Store   // Storage backend to keep track of token bucket states
	capacity        float64       // Maximum number of tokens in the bucket (burst capacity)
	refillRate      float64       // Tokens added to the bucket per second
	mutexes         sync.Map      // Map of mutexes for per-key synchronization
	cleanupTicker   *time.Ticker  // Ticker for periodic mutex cleanup
	cleanupStopCh   chan struct{} // Channel to stop the cleanup goroutine
	cleanupInterval time.Duration // Interval for mutex cleanup
}

// NewTokenBucketLimiter creates a new instance of TokenBucketLimiter.
//
// Parameters:
//   - store: Storage backend implementing the store.Store interface (e.g., RedisStore, MemoryStore)
//   - capacity: Maximum number of tokens in the bucket (burst capacity)
//   - refillRate: The rate at which tokens are added to the bucket (tokens per second)
//
// Returns:
//   - A pointer to a TokenBucketLimiter instance
func NewTokenBucketLimiter(store store.Store, capacity, refillRate float64) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		store:           store,
		capacity:        capacity,
		refillRate:      refillRate,
		mutexes:         sync.Map{},
		cleanupInterval: time.Minute * 5,
		cleanupStopCh:   make(chan struct{}),
	}
	go limiter.startMutexCleanup()
	return limiter
}

// getMutex returns the mutex associated with the key.
// If no mutex exists, it creates one.
func (l *TokenBucketLimiter) getMutex(key string) *keyMutex {
	mutexInterface, _ := l.mutexes.LoadOrStore(key, &keyMutex{
		mu:         &sync.Mutex{},
		lastAccess: time.Now(),
	})
	return mutexInterface.(*keyMutex)
}

// Allow checks whether a request associated with the given key is allowed under the rate limit.
// It refills tokens based on the elapsed time and consumes a token if available.
//
// Parameters:
//   - key: A unique identifier for the client (e.g., IP address, user ID)
//
// Returns:
//   - allowed: A boolean indicating whether the request is allowed (true) or should be rate-limited (false)
//   - err: An error if there was a problem accessing the storage backend
func (l *TokenBucketLimiter) Allow(key string) (bool, error) {
	km := l.getMutex(key)
	km.mu.Lock()
	defer km.mu.Unlock()
	km.lastAccess = time.Now()

	now := time.Now().UnixNano()

	// Retrieve the current token bucket state
	state, err := l.store.GetTokenBucket(key)
	if err != nil {
		return false, err
	}

	if state == nil {
		// Initialize a new token bucket state
		state = &store.TokenBucketState{
			Tokens:         l.capacity - 1, // Consume one token
			LastUpdateTime: now,
		}
		err = l.store.SetTokenBucket(key, state, time.Hour*24) // Set expiration as needed
		if err != nil {
			return false, err
		}
		return true, nil // Request is allowed
	}

	// Refill tokens based on the elapsed time
	elapsedTime := float64(now-state.LastUpdateTime) / float64(time.Second)
	refillTokens := elapsedTime * l.refillRate
	state.Tokens = min(state.Tokens+refillTokens, l.capacity)
	state.LastUpdateTime = now

	if state.Tokens >= 1 {
		// Consume a token
		state.Tokens -= 1
		err = l.store.SetTokenBucket(key, state, time.Hour*24)
		if err != nil {
			return false, err
		}
		return true, nil // Request is allowed
	}

	// Not enough tokens
	err = l.store.SetTokenBucket(key, state, time.Hour*24)
	if err != nil {
		return false, err
	}
	return false, nil // Rate limit exceeded
}

// startMutexCleanup runs a background goroutine to clean up unused mutexes.
func (l *TokenBucketLimiter) startMutexCleanup() {
	l.cleanupTicker = time.NewTicker(l.cleanupInterval)
	for {
		select {
		case <-l.cleanupTicker.C:
			now := time.Now()
			l.mutexes.Range(func(key, value interface{}) bool {
				km := value.(*keyMutex)
				km.mu.Lock()
				if now.Sub(km.lastAccess) > l.cleanupInterval*2 {
					km.mu.Unlock()
					l.mutexes.Delete(key)
				} else {
					km.mu.Unlock()
				}
				return true
			})
		case <-l.cleanupStopCh:
			l.cleanupTicker.Stop()
			return
		}
	}
}

// StopCleanup stops the mutex cleanup goroutine.
func (l *TokenBucketLimiter) StopCleanup() {
	close(l.cleanupStopCh)
}

// min returns the smaller of two float64 numbers.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
