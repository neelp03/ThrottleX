// store/store.go

package store

import "time"

// Store is an interface for storage backends used by rate limiters.
type Store interface {
	// Fixed Window methods
	Increment(key string, expiration time.Duration) (int64, error)

	// Sliding Window methods
	AddTimestamp(key string, timestamp int64, expiration time.Duration) error
	CountTimestamps(key string, start int64, end int64) (int64, error)

	// Token Bucket methods
	GetTokenBucket(key string) (*TokenBucketState, error)
	SetTokenBucket(key string, state *TokenBucketState, expiration time.Duration) error

	// Leaky Bucket methods
	GetLeakyBucket(key string) (*LeakyBucketState, error)
	SetLeakyBucket(key string, state *LeakyBucketState, expiration time.Duration) error
}

// TokenBucketState represents the state of a token bucket.
type TokenBucketState struct {
	Tokens         float64 // Current number of tokens in the bucket
	LastUpdateTime int64   // Unix timestamp in nanoseconds of the last update
}

// LeakyBucketState represents the state of a leaky bucket.
type LeakyBucketState struct {
	Queue        int       // Number of requests currently in the bucket
	LastLeakTime time.Time // Time when the bucket last leaked
}
