// ratelimiter/fixed_window.go

package ratelimiter

import (
	"errors"
	"strconv"
	"time"
	"regexp"

	"github.com/neelp03/throttlex/store"
)

// FixedWindowLimiter implements the fixed window rate limiting algorithm.
// It limits the number of requests allowed within a fixed time window.
// Once the limit is reached, all subsequent requests are denied until the window resets.
type FixedWindowLimiter struct {
	store  store.Store   // Storage backend to keep track of request counts
	limit  int           // Maximum number of requests allowed in the window
	window time.Duration // Duration of the fixed time window
}

// NewFixedWindowLimiter creates a new instance of FixedWindowLimiter.
//
// Parameters:
//   - store: Storage backend implementing the store.Store interface (e.g., RedisStore, MemoryStore)
//   - limit: Maximum number of requests allowed within the time window
//   - window: Duration of the fixed time window (e.g., time.Minute * 1)
//
// Returns:
//   - A pointer to a FixedWindowLimiter instance
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

// validKeyRegex is a compiled regular expression that matches valid keys.
// A valid key consists of alphanumeric characters, periods, underscores, and hyphens.
var validKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// Allow checks whether a request associated with the given key is allowed under the rate limit.
// It increments the count for the current window and determines if the request should be allowed.
//
// Parameters:
//   - key: A unique identifier for the client (e.g., IP address, user ID)
//
// Returns:
//   - allowed: A boolean indicating whether the request is allowed (true) or should be rate-limited (false)
//   - err: An error if there was a problem accessing the storage backend
func (l *FixedWindowLimiter) Allow(key string) (bool, error) {
	// Input validation

	// Check for empty key
	if key == "" {
		return false, errors.New("invalid key: key cannot be empty")
	}

	// Check for overly long key (256 characters max)
	if len(key) > 256 {
		return false, errors.New("invalid key: key length exceeds maximum allowed length")
	}

	// Check for valid key format (alphanumeric, ".", "_", "-")
	if !validKeyRegex.MatchString(key) {
		return false, errors.New("invalid key: key contains invalid characters")
	}

	// Proceed with rate limiting if input validation passes
	windowKey := l.getWindowKey(key)
	count, err := l.store.Increment(windowKey, l.window)
	if err != nil {
		return false, err
	}

	if count > int64(l.limit) {
		return false, nil // Rate limit exceeded
	}
	return true, nil // Request is allowed
}

// getWindowKey generates a unique key for the current time window and client key.
// This ensures counts are tracked separately for each client and time window.
//
// Parameters:
//   - key: The unique identifier for the client
//
// Returns:
//   - A string representing the combined key for the client and current window
func (l *FixedWindowLimiter) getWindowKey(key string) string {
	// Calculate the current window number based on the time
	windowNumber := time.Now().Unix() / int64(l.window.Seconds())

	// Combine the client key with the window number to form a unique key
	return key + ":" + strconv.FormatInt(windowNumber, 10)
}
