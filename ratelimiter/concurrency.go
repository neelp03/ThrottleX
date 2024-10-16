// ratelimiter/concurrency.go

package ratelimiter

import (
	"errors"
)

// ConcurrencyLimiter limits the number of concurrent requests.
type ConcurrencyLimiter struct {
	limit     int
	semaphore chan struct{}
}

// NewConcurrencyLimiter creates a new ConcurrencyLimiter.
func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	if limit <= 0 {
		limit = 1
	}
	return &ConcurrencyLimiter{
		limit:     limit,
		semaphore: make(chan struct{}, limit),
	}
}

// Allow checks if a request is allowed based on the concurrency limit.
func (l *ConcurrencyLimiter) Allow(key string) (bool, error) {
	// Input validation
	if key == "" {
		return false, errors.New("invalid key: key cannot be empty")
	}
	if len(key) > 256 {
		return false, errors.New("invalid key: key length exceeds maximum allowed length")
	}
	if !validKeyRegex.MatchString(key) {
		return false, errors.New("invalid key: key contains invalid characters")
	}

	select {
	case l.semaphore <- struct{}{}:
		return true, nil
	default:
		return false, nil
	}
}

// Release should be called when the request is finished processing.
func (l *ConcurrencyLimiter) Release() error {
	select {
	case <-l.semaphore:
		return nil
	default:
		return errors.New("semaphore underflow")
	}
}
