package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/neelp03/throttlex/ratelimiter"
	"github.com/neelp03/throttlex/store"
)

func main() {
	// Initialize the store (using MemoryStore for simplicity)
	memStore := store.NewMemoryStore()

	// Create a FixedWindowLimiter allowing 100 requests per minute
	limiter := ratelimiter.NewFixedWindowLimiter(memStore, 100, time.Minute)

	// Create a rate-limiting middleware
	rateLimitMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr // Use client IP as the key

			allowed, err := limiter.Allow(key)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !allowed {
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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
