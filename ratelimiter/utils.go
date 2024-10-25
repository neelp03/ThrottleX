package ratelimiter

import (
	"errors"
	"regexp"
)

// validKeyRegex is a compiled regular expression that matches valid keys.
var validKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// validateKey checks if the provided key is valid.
func validateKey(key string) error {
	if key == "" {
		return errors.New("invalid key: key cannot be empty")
	}
	if len(key) > 256 {
		return errors.New("invalid key: key length exceeds maximum allowed length")
	}
	if !validKeyRegex.MatchString(key) {
		return errors.New("invalid key: key contains invalid characters")
	}
	return nil
}

// min returns the smaller of two float64 numbers.
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
