// ratelimiter/ratelimiter.go
package ratelimiter

// RateLimiter is an interface that defines the contract for rate limiting algorithms.
// It provides a method to determine whether a request identified by a unique key is allowed
// based on the implemented rate limiting strategy.
type RateLimiter interface {
	// Allow checks if a request associated with the given key is allowed to proceed.
	// It returns true if the request is allowed, false otherwise.
	// An error is returned if there was an issue checking the rate limit.
	Allow(key string) (bool, error)
}
