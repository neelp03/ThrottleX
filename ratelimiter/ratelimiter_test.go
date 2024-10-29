package ratelimiter

import (
	"testing"
	"time"

	"github.com/neelp03/throttlex/store"
)

func TestNewRateLimiter(t *testing.T) {
	memStore := store.NewMemoryStore()

	tests := []struct {
		name        string
		config      LimiterConfig
		expectError bool
	}{
		{
			name: "FixedWindowPolicy",
			config: LimiterConfig{
				Policy:   FixedWindowPolicy,
				Store:    memStore,
				Limit:    5,
				Interval: time.Second,
			},
			expectError: false,
		},
		{
			name: "SlidingWindowPolicy",
			config: LimiterConfig{
				Policy:   SlidingWindowPolicy,
				Store:    memStore,
				Limit:    5,
				Interval: time.Second,
			},
			expectError: false,
		},
		{
			name: "TokenBucketPolicy",
			config: LimiterConfig{
				Policy:     TokenBucketPolicy,
				Store:      memStore,
				Capacity:   10,
				RefillRate: 1,
			},
			expectError: false,
		},
		{
			name: "LeakyBucketPolicy",
			config: LimiterConfig{
				Policy:   LeakyBucketPolicy,
				Store:    memStore,
				Capacity: 10,
				LeakRate: 1,
			},
			expectError: false,
		},
		{
			name: "ConcurrencyPolicy",
			config: LimiterConfig{
				Policy:      ConcurrencyPolicy,
				Store:       memStore,
				Concurrency: 3,
			},
			expectError: false,
		},
		{
			name: "UnknownPolicy",
			config: LimiterConfig{
				Policy: "UnknownPolicy",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter, err := NewRateLimiter(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check if the correct type of limiter is created for each policy
			switch tt.config.Policy {
			case FixedWindowPolicy:
				if _, ok := limiter.(*FixedWindowLimiter); !ok {
					t.Errorf("Expected FixedWindowLimiter, got %T", limiter)
				}
			case SlidingWindowPolicy:
				if _, ok := limiter.(*SlidingWindowLimiter); !ok {
					t.Errorf("Expected SlidingWindowLimiter, got %T", limiter)
				}
			case TokenBucketPolicy:
				if _, ok := limiter.(*TokenBucketLimiter); !ok {
					t.Errorf("Expected TokenBucketLimiter, got %T", limiter)
				}
			case LeakyBucketPolicy:
				if _, ok := limiter.(*LeakyBucketLimiter); !ok {
					t.Errorf("Expected LeakyBucketLimiter, got %T", limiter)
				}
			case ConcurrencyPolicy:
				if _, ok := limiter.(*ConcurrencyLimiter); !ok {
					t.Errorf("Expected ConcurrencyLimiter, got %T", limiter)
				}
			}
		})
	}
}
