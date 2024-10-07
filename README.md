  # **Throttlex**

[![CI](https://github.com/neelp03/throttlex/actions/workflows/ci.yml/badge.svg)](https://github.com/neelp03/throttlex/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/neelp03/throttlex/branch/main/graph/badge.svg)](https://codecov.io/gh/neelp03/throttlex)
[![Go Report Card](https://goreportcard.com/badge/github.com/neelp03/throttlex?v=1)](https://goreportcard.com/report/github.com/neelp03/throttlex)
[![GoDoc](https://godoc.org/github.com/neelp03/throttlex?status.svg)](https://godoc.org/github.com/neelp03/throttlex)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**Throttlex** is a Go package that provides flexible and efficient rate limiting for your applications. It supports multiple rate-limiting algorithms and storage backends, making it suitable for a variety of use cases, including REST APIs, gRPC services, and GraphQL APIs.

---

## **Table of Contents**

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Integrate with REST API](#integrate-with-rest-api)
  - [Integrate with gRPC Service](#integrate-with-grpc-service)
  - [Integrate with GraphQL API](#integrate-with-graphql-api)
- [Rate Limiting Algorithms](#rate-limiting-algorithms)
- [Storage Backends](#storage-backends)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

---

## **Features**

- **Multiple Algorithms**:
  - Fixed Window
  - Sliding Window
  - Token Bucket

- **Pluggable Storage Backends**:
  - In-Memory Store
  - Redis Store

- **Thread-Safe and Efficient**:
  - Designed for high concurrency and low latency.
  - Includes mutex cleanup to prevent memory leaks.

- **Easy Integration**:
  - Middleware examples for REST APIs.
  - Interceptor examples for gRPC services.
  - Middleware examples for GraphQL APIs.

- **Highly Configurable**:
  - Customize limits, window sizes, and keys.
  - Support for dynamic configurations.

---

## **Installation**

To install Throttlex, use `go get`:

```bash
go get -u github.com/neelp03/throttlex
```

---

## **Usage**

Import the package into your Go project:

```go
import (
    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)
```

### **Integrate with REST API**

Here's how to integrate Throttlex with a REST API using the `net/http` package:

```go
package main

import (
    "fmt"
    "net"
    "net/http"
    "time"

    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)

func main() {
    // Initialize the store
    memStore := store.NewMemoryStore()

    // Create a FixedWindowLimiter
    limit := 100
    window := time.Minute
    limiter, err := ratelimiter.NewFixedWindowLimiter(memStore, limit, window)
    if err != nil {
        fmt.Printf("Failed to create rate limiter: %v\n", err)
        return
    }

    // Rate-limiting middleware
    rateLimitMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract client IP
            ip, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            key := ip

            allowed, err := limiter.Allow(key)
            if err != nil {
                http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
                return
            }
            if !allowed {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }

    // Define your handler
    helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, World!")
    })

    // Apply middleware
    http.Handle("/", rateLimitMiddleware(helloHandler))

    // Start the server
    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

### **Integrate with gRPC Service**

Here's how to use Throttlex with a gRPC service:

```go
package main

import (
    "context"
    "fmt"
    "net"

    "google.golang.org/grpc"
    "google.golang.org/grpc/peer"

    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)

func main() {
    // Initialize the store
    memStore := store.NewMemoryStore()

    // Create a TokenBucketLimiter
    capacity := 20
    refillRate := 10
    limiter, err := ratelimiter.NewTokenBucketLimiter(memStore, float64(capacity), float64(refillRate))
    if err != nil {
        fmt.Printf("Failed to create rate limiter: %v\n", err)
        return
    }

    // Rate-limiting interceptor
    rateLimitInterceptor := func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        // Extract client IP or other identifier
        p, ok := peer.FromContext(ctx)
        var key string
        if ok {
            key = p.Addr.String()
        } else {
            key = "unknown"
        }

        allowed, err := limiter.Allow(key)
        if err != nil {
            return nil, grpc.Errorf(grpc.Code(grpc.Internal), "Internal Server Error")
        }
        if !allowed {
            return nil, grpc.Errorf(grpc.Code(grpc.ResourceExhausted), "Too Many Requests")
        }

        return handler(ctx, req)
    }

    // Create a gRPC server with the interceptor
    serverOptions := []grpc.ServerOption{
        grpc.UnaryInterceptor(rateLimitInterceptor),
    }
    grpcServer := grpc.NewServer(serverOptions...)

    // Register your services
    // pb.RegisterYourServiceServer(grpcServer, &YourServiceServer{})

    // Start the server
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        fmt.Printf("Failed to listen: %v\n", err)
        return
    }
    fmt.Println("gRPC server is running on port 50051")
    grpcServer.Serve(lis)
}
```

### **Integrate with GraphQL API**

Integrate Throttlex with a GraphQL API using `graphql-go`:

```go
package main

import (
    "fmt"
    "net"
    "net/http"
    "time"

    "github.com/graphql-go/graphql"
    "github.com/graphql-go/handler"

    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)

func main() {
    // Define GraphQL schema
    fields := graphql.Fields{
        "hello": &graphql.Field{
            Type: graphql.String,
            Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                return "Hello, World!", nil
            },
        },
    }
    rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
    schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
    schema, err := graphql.NewSchema(schemaConfig)
    if err != nil {
        fmt.Printf("Failed to create schema: %v\n", err)
        return
    }

    // Initialize the store and rate limiter
    memStore := store.NewMemoryStore()
    limiter, err := ratelimiter.NewSlidingWindowLimiter(memStore, 100, time.Minute)
    if err != nil {
        fmt.Printf("Failed to create rate limiter: %v\n", err)
        return
    }

    // Rate-limiting middleware
    rateLimitMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            key := ip

            allowed, err := limiter.Allow(key)
            if err != nil {
                http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
                return
            }
            if !allowed {
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }

            next.ServeHTTP(w, r)
        })
    }

    // Create GraphQL handler
    h := handler.New(&handler.Config{
        Schema: &schema,
        Pretty: true,
    })

    // Apply middleware
    http.Handle("/graphql", rateLimitMiddleware(h))

    // Start the server
    fmt.Println("GraphQL server is running on http://localhost:8080/graphql")
    http.ListenAndServe(":8080", nil)
}
```

---

## **Rate Limiting Algorithms**

### **1. Fixed Window Limiter**

Counts the number of requests in fixed time intervals. Simple but can be prone to spikes at window boundaries.

**Constructor:**

```go
limiter, err := ratelimiter.NewFixedWindowLimiter(store, limit, window)
```

- `store`: Storage backend (`MemoryStore`, `RedisStore`).
- `limit`: Maximum number of requests allowed in the window.
- `window`: Duration of the window (e.g., `time.Minute`).

### **2. Sliding Window Limiter**

Provides a smoother rate limit by using a sliding time window, counting requests over the last N seconds.

**Constructor:**

```go
limiter, err := ratelimiter.NewSlidingWindowLimiter(store, limit, window)
```

### **3. Token Bucket Limiter**

Allows bursts of traffic by storing tokens that replenish at a fixed rate.

**Constructor:**

```go
limiter, err := ratelimiter.NewTokenBucketLimiter(store, capacity, refillRate)
```

- `capacity`: Maximum number of tokens in the bucket.
- `refillRate`: Number of tokens added per second.

---

## **Storage Backends**

### **MemoryStore**

An in-memory store suitable for single-instance applications.

**Usage:**

```go
memStore := store.NewMemoryStore()
```

### **RedisStore**

A Redis-based store suitable for distributed systems.

**Usage:**

```go
import "github.com/go-redis/redis/v8"

redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
redisStore := store.NewRedisStore(redisClient)
```

- `redisClient`: An instance of `*redis.Client` from the `go-redis` library.

---

## **Configuration**

Customize your rate limiter by adjusting the following parameters:

- **Limit**: Maximum number of requests.
- **Window Size**: Time duration for the rate limit window.
- **Capacity and Refill Rate**: For the Token Bucket algorithm.
- **Keys**: Unique identifiers for clients (e.g., IP address, user ID).

---

## **Examples**

Explore the `examples` directory for full example applications:

- **REST API**: [`examples/rest-api`](examples/rest-api)
- **gRPC Service**: [`examples/grpc-api`](examples/grpc-api)
- **GraphQL API**: [`examples/graphql-api`](examples/graphql-api)

---

## **Contributing**

Contributions are welcome! Please follow these steps:

1. **Fork the Repository**: Click on the "Fork" button at the top.
2. **Clone Your Fork**:

   ```bash
   git clone https://github.com/neelp03/throttlex.git
   cd throttlex
   ```

3. **Create a Branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make Changes**: Implement your feature or fix.
5. **Run Tests**:

   ```bash
   go test -race -v ./...
   ```

6. **Commit and Push**:

   ```bash
   git add .
   git commit -m "Add your feature"
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**: Open a pull request against the `main` branch.

---

## **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## **Acknowledgments**

- Inspired by the need for flexible and efficient rate limiting in Go applications.
- Thanks to the Go community for their invaluable contributions.

---

## **Contact**

For questions or support, please open an issue on the [GitHub repository](https://github.com/neelp03/throttlex/issues).
