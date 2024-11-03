package ratelimiter

import (
	"errors"
	"sync"
	"time"

	"github.com/neelp03/throttlex/store"
)

// SlidingWindowLimiter implements the sliding window rate-limiting algorithm.
type SlidingWindowLimiter struct {
	store           store.Store
	limit           int
	window          time.Duration
	mutexes         sync.Map
	cleanupTicker   *time.Ticker
	cleanupStopCh   chan struct{}
	cleanupInterval time.Duration
}

// NewSlidingWindowLimiter creates a new instance of SlidingWindowLimiter.
func NewSlidingWindowLimiter(store store.Store, limit int, window time.Duration) (*SlidingWindowLimiter, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}
	if window <= 0 {
		return nil, errors.New("window duration must be greater than zero")
	}
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	limiter := &SlidingWindowLimiter{
		store:           store,
		limit:           limit,
		window:          window,
		mutexes:         sync.Map{},
		cleanupInterval: time.Minute * 5,
		cleanupStopCh:   make(chan struct{}),
	}
	go limiter.startMutexCleanup()
	return limiter, nil
}

// getMutex returns the mutex associated with the key.
func (l *SlidingWindowLimiter) getMutex(key string) *keyMutex {
	mutexInterface, _ := l.mutexes.LoadOrStore(key, &keyMutex{
		mu:         &sync.Mutex{},
		lastAccess: time.Now(),
	})
	return mutexInterface.(*keyMutex)
}

// Allow checks whether a request associated with the given key is allowed.
func (l *SlidingWindowLimiter) Allow(key string) (bool, error) {
	if err := validateKey(key); err != nil {
		return false, err
	}

	km := l.getMutex(key)
	km.mu.Lock()
	defer km.mu.Unlock()
	km.lastAccess = time.Now()

	now := time.Now().UnixNano()
	windowStart := now - l.window.Nanoseconds()

	count, err := l.store.CountTimestamps(key, windowStart, now)
	if err != nil {
		return false, err
	}

	allowed := count < int64(l.limit)
	if !allowed {
		return false, nil
	}

	err = l.store.AddTimestamp(key, now, l.window)
	if err != nil {
		return false, err
	}

	return true, nil
}

// startMutexCleanup runs a background goroutine to clean up unused mutexes.
func (l *SlidingWindowLimiter) startMutexCleanup() {
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
func (l *SlidingWindowLimiter) StopCleanup() {
	close(l.cleanupStopCh)
}
