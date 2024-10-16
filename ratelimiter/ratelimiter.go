// ratelimiter/ratelimiter.go

package ratelimiter

import (
	"fmt"
	"time"

	"github.com/neelp03/throttlex/store"
)

// RateLimiter is an interface that defines the contract for rate limiting algorithms.
// It provides a method to determine whether a request identified by a unique key is allowed
// based on the implemented rate limiting strategy.
type RateLimiter interface {
	// Allow checks if a request associated with the given key is allowed to proceed.
	// It returns true if the request is allowed, false otherwise.
	// An error is returned if there was an issue checking the rate limit.
	Allow(key string) (bool, error)
}

// PolicyType represents the type of rate-limiting policy.
type PolicyType string

const (
	FixedWindow   PolicyType = "FixedWindow"
	SlidingWindow PolicyType = "SlidingWindow"
	TokenBucket   PolicyType = "TokenBucket"
	LeakyBucket   PolicyType = "LeakyBucket"
	Concurrency   PolicyType = "Concurrency"
)

// LimiterConfig holds configuration for a rate limiter.
type LimiterConfig struct {
	Policy     PolicyType
	Store      store.Store
	Limit      int           // General limit parameter
	Interval   time.Duration // General interval parameter
	Capacity   float64       // For TokenBucket and LeakyBucket
	RefillRate float64       // For TokenBucket
	LeakRate   time.Duration // For LeakyBucket
}

// NewRateLimiter is a factory function that creates a RateLimiter based on the specified policy.
func NewRateLimiter(config LimiterConfig) (RateLimiter, error) {
	switch config.Policy {
	case FixedWindow:
		return NewFixedWindowLimiter(config.Store, config.Limit, config.Interval)
	case SlidingWindow:
		return NewSlidingWindowLimiter(config.Store, config.Limit, config.Interval)
	case TokenBucket:
		return NewTokenBucketLimiter(config.Store, config.Capacity, config.RefillRate)
	case LeakyBucket:
		return NewLeakyBucketLimiter(config.Store, int(config.Capacity), config.LeakRate)
	case Concurrency:
		return NewConcurrencyLimiter(config.Limit), nil
	default:
		return nil, fmt.Errorf("unknown rate limiting policy: %s", config.Policy)
	}
}
