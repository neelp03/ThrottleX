package ratelimiter

import (
	"errors"
	"strconv"
	"time"

	"github.com/neelp03/throttlex/store"
)

// FixedWindowLimiter implements the fixed window rate limiting algorithm.
type FixedWindowLimiter struct {
	store  store.Store   // Storage backend to keep track of request counts
	limit  int           // Maximum number of requests allowed in the window
	window time.Duration // Duration of the fixed time window
}

// NewFixedWindowLimiter creates a new instance of FixedWindowLimiter.
func NewFixedWindowLimiter(store store.Store, limit int, window time.Duration) (*FixedWindowLimiter, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be greater than zero")
	}
	if window <= 0 {
		return nil, errors.New("window duration must be greater than zero")
	}
	if store == nil {
		return nil, errors.New("store cannot be nil")
	}

	return &FixedWindowLimiter{
		store:  store,
		limit:  limit,
		window: window,
	}, nil
}

// Allow checks whether a request associated with the given key is allowed under the rate limit.
func (l *FixedWindowLimiter) Allow(key string) (bool, error) {
	// Input validation
	if err := validateKey(key); err != nil {
		return false, err
	}

	// Proceed with rate limiting if input validation passes
	windowKey := l.getWindowKey(key)
	count, err := l.store.Increment(windowKey, 1, l.window) // Added '1' as delta parameter
	if err != nil {
		return false, err
	}

	if count > int64(l.limit) {
		return false, nil // Rate limit exceeded
	}
	return true, nil // Request is allowed
}

// getWindowKey generates a unique key for the current time window and client key.
func (l *FixedWindowLimiter) getWindowKey(key string) string {
	// Calculate the current window number based on the time
	windowNumber := time.Now().Unix() / int64(l.window.Seconds())

	// Combine the client key with the window number to form a unique key
	return key + ":" + strconv.FormatInt(windowNumber, 10)
}
