// store/store.go

package store

import "time"

// Store is an interface for storage backends used by rate limiters.
type Store interface {
	// Increment increments the counter and sets expiration.
	Increment(key string, expiration time.Duration) (int64, error)

	// AddTimestamp adds a timestamp to the list associated with the key.
	AddTimestamp(key string, timestamp int64, expiration time.Duration) error

	// CountTimestamps counts the number of timestamps within the given range [start, end].
	CountTimestamps(key string, start int64, end int64) (int64, error)

	// GetTokenBucket retrieves the current state of the token bucket.
	GetTokenBucket(key string) (*TokenBucketState, error)

	// SetTokenBucket updates the state of the token bucket.
	SetTokenBucket(key string, state *TokenBucketState, expiration time.Duration) error
}

// TokenBucketState represents the state of a token bucket.
// It includes the number of tokens currently available and the last update time.
type TokenBucketState struct {
	Tokens         float64 // Current number of tokens in the bucket
	LastUpdateTime int64   // Unix timestamp in nanoseconds of the last update
}
