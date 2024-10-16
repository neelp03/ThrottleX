// ratelimiter/leaky_bucket.go

package ratelimiter

import (
	"errors"
	"sync"
	"time"

	"github.com/neelp03/throttlex/store"
)

// LeakyBucketLimiter implements the leaky bucket rate-limiting algorithm.
type LeakyBucketLimiter struct {
	store           store.Store
	capacity        int
	leakRate        time.Duration
	mutexes         sync.Map
	cleanupTicker   *time.Ticker
	cleanupStopCh   chan struct{}
	cleanupInterval time.Duration
}

// NewLeakyBucketLimiter creates a new LeakyBucketLimiter.
func NewLeakyBucketLimiter(store store.Store, capacity int, leakRate time.Duration) (*LeakyBucketLimiter, error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be greater than zero")
	}
	if leakRate <= 0 {
		return nil, errors.New("leakRate must be greater than zero")
	}
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	limiter := &LeakyBucketLimiter{
		store:           store,
		capacity:        capacity,
		leakRate:        leakRate,
		cleanupInterval: time.Minute * 5,
		cleanupStopCh:   make(chan struct{}),
	}
	go limiter.startMutexCleanup()
	return limiter, nil
}

// Allow checks whether a request is allowed under the leaky bucket rate limit.
func (l *LeakyBucketLimiter) Allow(key string) (bool, error) {
	// Input validation
	if err := validateKey(key); err != nil {
		return false, err
	}

	km := l.getMutex(key)
	km.mu.Lock()
	defer km.mu.Unlock()
	km.lastAccess = time.Now()

	// Get the current state from the store
	state, err := l.store.GetLeakyBucket(key)
	if err != nil {
		return false, err
	}

	now := time.Now()
	if state == nil {
		// Initialize new bucket
		state = &store.LeakyBucketState{
			LastLeakTime: now,
			Queue:        0,
		}
	}

	// Calculate how many requests have leaked since last check
	elapsed := now.Sub(state.LastLeakTime)
	leaked := int(elapsed / l.leakRate)
	if leaked > 0 {
		state.Queue = max(0, state.Queue-leaked)
		state.LastLeakTime = state.LastLeakTime.Add(time.Duration(leaked) * l.leakRate)
	}

	// Check if there's capacity in the bucket
	if state.Queue < l.capacity {
		// Add request to the bucket
		state.Queue++
		err = l.store.SetLeakyBucket(key, state, time.Hour*24)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Bucket is full
	err = l.store.SetLeakyBucket(key, state, time.Hour*24)
	if err != nil {
		return false, err
	}
	return false, nil
}

// getMutex returns the mutex associated with the key.
func (l *LeakyBucketLimiter) getMutex(key string) *keyMutex {
	mutexInterface, _ := l.mutexes.LoadOrStore(key, &keyMutex{
		mu:         &sync.Mutex{},
		lastAccess: time.Now(),
	})
	return mutexInterface.(*keyMutex)
}

// startMutexCleanup runs a background goroutine to clean up unused mutexes.
func (l *LeakyBucketLimiter) startMutexCleanup() {
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
func (l *LeakyBucketLimiter) StopCleanup() {
	close(l.cleanupStopCh)
}
