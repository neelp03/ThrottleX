package ratelimiter

import (
	"errors"
	"sync"
	"time"

	"github.com/neelp03/throttlex/store"
)

// ConcurrencyLimiter limits the number of concurrent requests per key.
type ConcurrencyLimiter struct {
	store           store.Store   // Storage backend to keep track of concurrency counts
	maxConcurrent   int64         // Maximum number of concurrent requests
	mutexes         sync.Map      // Map of mutexes for per-key synchronization
	cleanupTicker   *time.Ticker  // Ticker for periodic mutex cleanup
	cleanupStopCh   chan struct{} // Channel to stop the cleanup goroutine
	cleanupInterval time.Duration // Interval for mutex cleanup
}

// NewConcurrencyLimiter creates a new ConcurrencyLimiter.
func NewConcurrencyLimiter(store store.Store, maxConcurrent int64) (*ConcurrencyLimiter, error) {
	if maxConcurrent <= 0 {
		return nil, errors.New("maxConcurrent must be greater than zero")
	}
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	limiter := &ConcurrencyLimiter{
		store:           store,
		maxConcurrent:   maxConcurrent,
		mutexes:         sync.Map{},
		cleanupInterval: time.Minute * 5,
		cleanupStopCh:   make(chan struct{}),
	}
	go limiter.startMutexCleanup()
	return limiter, nil
}

// getMutex returns the mutex associated with the key.
func (cl *ConcurrencyLimiter) getMutex(key string) *keyMutex {
	mutexInterface, _ := cl.mutexes.LoadOrStore(key, &keyMutex{
		mu:         &sync.Mutex{},
		lastAccess: time.Now(),
	})
	return mutexInterface.(*keyMutex)
}

// Allow tries to acquire a slot for processing.
func (cl *ConcurrencyLimiter) Allow(key string) (bool, error) {
	// Input validation
	if err := validateKey(key); err != nil {
		return false, err
	}

	km := cl.getMutex(key)
	km.mu.Lock()
	defer km.mu.Unlock()
	km.lastAccess = time.Now()

	count, err := cl.store.Increment(key, 1, time.Hour*24)
	if err != nil {
		return false, err
	}

	if count > cl.maxConcurrent {
		// Exceeded limit, decrement the count
		_, err = cl.store.Increment(key, -1, time.Hour*24)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

// Release releases a slot after processing.
func (cl *ConcurrencyLimiter) Release(key string) error {
	km := cl.getMutex(key)
	km.mu.Lock()
	defer km.mu.Unlock()
	km.lastAccess = time.Now()

	_, err := cl.store.Increment(key, -1, time.Hour*24)
	return err
}

// startMutexCleanup runs a background goroutine to clean up unused mutexes.
func (cl *ConcurrencyLimiter) startMutexCleanup() {
	cl.cleanupTicker = time.NewTicker(cl.cleanupInterval)
	for {
		select {
		case <-cl.cleanupTicker.C:
			now := time.Now()
			cl.mutexes.Range(func(key, value interface{}) bool {
				km := value.(*keyMutex)
				km.mu.Lock()
				if now.Sub(km.lastAccess) > cl.cleanupInterval*2 {
					km.mu.Unlock()
					cl.mutexes.Delete(key)
				} else {
					km.mu.Unlock()
				}
				return true
			})
		case <-cl.cleanupStopCh:
			cl.cleanupTicker.Stop()
			return
		}
	}
}

// StopCleanup stops the mutex cleanup goroutine.
func (cl *ConcurrencyLimiter) StopCleanup() {
	close(cl.cleanupStopCh)
}
