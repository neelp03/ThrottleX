package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/neelp03/throttlex/ratelimiter"
	"github.com/neelp03/throttlex/store"
)

func main() {
	// Initialize the store
	memStore := store.NewMemoryStore()

	// Create a FixedWindowLimiter with proper parameter validation
	limit := 100
	window := time.Minute
	limiter, err := ratelimiter.NewFixedWindowLimiter(memStore, limit, window)
	if err != nil {
		fmt.Printf("Failed to create rate limiter: %v\n", err)
		return
	}

	// Rate-limiting middleware with enhanced edge case handling
	rateLimitMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key, err := getClientIP(r)
			if err != nil {
				http.Error(w, "Unable to determine client IP", http.StatusBadRequest)
				return
			}

			allowed, err := limiter.Allow(key)
			if err != nil {
				http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
				return
			}
			if !allowed {
				retryAfter := window.Seconds()
				w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Define your HTTP handler
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})

	// Wrap your handler with the middleware
	http.Handle("/", rateLimitMiddleware(helloHandler))

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

// getClientIP extracts the client IP address, considering proxies
func getClientIP(r *http.Request) (string, error) {
	// Check X-Forwarded-For header (assuming trusted proxy)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		return ip, nil
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("failed to parse IP from RemoteAddr: %v", err)
	}
	return ip, nil
}
