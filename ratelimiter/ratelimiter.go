package ratelimiter

import (
	"fmt"
	"time"

	"github.com/neelp03/throttlex/store"
)

// RateLimiter is an interface that defines the contract for rate limiting algorithms.
type RateLimiter interface {
	// Allow checks if a request associated with the given key is allowed to proceed.
	Allow(key string) (bool, error)
}

// PolicyType represents the type of rate-limiting policy.
type PolicyType string

const (
	FixedWindowPolicy   PolicyType = "FixedWindow"
	SlidingWindowPolicy PolicyType = "SlidingWindow"
	TokenBucketPolicy   PolicyType = "TokenBucket"
	LeakyBucketPolicy   PolicyType = "LeakyBucket"
	ConcurrencyPolicy   PolicyType = "Concurrency"
)

// LimiterConfig holds configuration for a rate limiter.
type LimiterConfig struct {
	Policy      PolicyType
	Store       store.Store
	Limit       int           // General limit parameter
	Interval    time.Duration // General interval parameter
	Capacity    float64       // For TokenBucket and LeakyBucket
	RefillRate  float64       // For TokenBucket
	LeakRate    float64       // For LeakyBucket
	Concurrency int64         // For ConcurrencyLimiter
}

// NewRateLimiter is a factory function that creates a RateLimiter based on the specified policy.
func NewRateLimiter(config LimiterConfig) (RateLimiter, error) {
	switch config.Policy {
	case FixedWindowPolicy:
		return NewFixedWindowLimiter(config.Store, config.Limit, config.Interval)
	case SlidingWindowPolicy:
		return NewSlidingWindowLimiter(config.Store, config.Limit, config.Interval)
	case TokenBucketPolicy:
		return NewTokenBucketLimiter(config.Store, config.Capacity, config.RefillRate)
	case LeakyBucketPolicy:
		return NewLeakyBucketLimiter(config.Store, int(config.Capacity), config.LeakRate)
	case ConcurrencyPolicy:
		return NewConcurrencyLimiter(config.Store, config.Concurrency)
	default:
		return nil, fmt.Errorf("unknown rate limiting policy: %s", config.Policy)
	}
}
