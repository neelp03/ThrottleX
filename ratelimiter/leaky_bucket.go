package ratelimiter

import (
	"errors"
	"sync"
	"time"

	"github.com/neelp03/throttlex/store"
)

// LeakyBucketLimiter implements the leaky bucket rate-limiting algorithm.
type LeakyBucketLimiter struct {
	store           store.Store   // Storage backend to keep track of leaky bucket states
	capacity        int           // Maximum capacity of the bucket
	leakRate        float64       // Leak rate per second
	mutexes         sync.Map      // Map of mutexes for per-key synchronization
	cleanupTicker   *time.Ticker  // Ticker for periodic mutex cleanup
	cleanupStopCh   chan struct{} // Channel to stop the cleanup goroutine
	cleanupInterval time.Duration // Interval for mutex cleanup
}

// NewLeakyBucketLimiter creates a new instance of LeakyBucketLimiter.
func NewLeakyBucketLimiter(store store.Store, capacity int, leakRate float64) (*LeakyBucketLimiter, error) {
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
		mutexes:         sync.Map{},
		cleanupInterval: time.Minute * 5,
		cleanupStopCh:   make(chan struct{}),
	}
	go limiter.startMutexCleanup()
	return limiter, nil
}

// getMutex returns the mutex associated with the key.
func (l *LeakyBucketLimiter) getMutex(key string) *keyMutex {
	mutexInterface, _ := l.mutexes.LoadOrStore(key, &keyMutex{
		mu:         &sync.Mutex{},
		lastAccess: time.Now(),
	})
	return mutexInterface.(*keyMutex)
}

// Allow checks whether a request associated with the given key is allowed under the rate limit.
func (l *LeakyBucketLimiter) Allow(key string) (bool, error) {
	// Input validation
	if err := validateKey(key); err != nil {
		return false, err
	}

	km := l.getMutex(key)
	km.mu.Lock()
	defer km.mu.Unlock()
	km.lastAccess = time.Now()

	state, err := l.store.GetLeakyBucket(key)
	if err != nil {
		return false, err
	}

	now := time.Now()

	if state == nil {
		// Initialize state
		state = &store.LeakyBucketState{
			Queue:        0,
			LastLeakTime: now,
		}
	} else {
		// Leak tokens based on elapsed time
		elapsed := now.Sub(state.LastLeakTime).Seconds()
		leaked := int(elapsed * l.leakRate)
		if leaked > 0 {
			state.Queue -= leaked
			if state.Queue < 0 {
				state.Queue = 0
			}
			// Update LastLeakTime
			state.LastLeakTime = state.LastLeakTime.Add(time.Duration(float64(leaked)/l.leakRate) * time.Second)
		}
	}

	if state.Queue < l.capacity {
		state.Queue++
		err = l.store.SetLeakyBucket(key, state, time.Hour*24)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		// Update the state even if not allowed
		err = l.store.SetLeakyBucket(key, state, time.Hour*24)
		if err != nil {
			return false, err
		}
		return false, nil
	}
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
