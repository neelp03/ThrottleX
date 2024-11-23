# ThrottleX: Scalable Rate Limiting for Go APIs

[![CI](https://github.com/neelp03/throttlex/actions/workflows/ci.yml/badge.svg)](https://github.com/neelp03/throttlex/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/neelp03/throttlex/branch/main/graph/badge.svg)](https://codecov.io/gh/neelp03/throttlex)
[![Go Report Card](https://goreportcard.com/badge/github.com/neelp03/throttlex?v=1)](https://goreportcard.com/report/github.com/neelp03/throttlex)
[![GoDoc](https://godoc.org/github.com/neelp03/throttlex?status.svg)](https://godoc.org/github.com/neelp03/throttlex)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

---

## **Overview**

ThrottleX is an advanced, flexible rate-limiting library designed to handle high traffic loads across distributed environments. Built in Go, ThrottleX combines multiple rate-limiting algorithms with concurrency control to ensure efficient API request management and system stability.

For more detailed setup instructions and examples, visit the **[ThrottleX Wiki](https://github.com/neelp03/ThrottleX/wiki)**.

---

## **Key Features**

- **Enhanced Rate Limiting Algorithms**:
  - Fixed Window
  - Sliding Window
  - Token Bucket
  - Leaky Bucket with Concurrency Limiting (new in `v1.0.0-rc2`)

- **Multiple Storage Options**:
  - In-Memory Store
  - Redis Store, optimized for distributed setups

- **Optimized for Performance**:
  - Concurrency control with goroutines and mutex management.
  - Customizable request limits, intervals, and policies.
  - Efficient memory management and optimized Redis configuration.

- **Real-Time Monitoring**:
  - Integrated with Prometheus for metrics collection and Grafana for visualization.

---

## **Installation**

Install ThrottleX via `go get`:

```bash
go get -u github.com/neelp03/throttlex
```

See the **[Installation and Setup Wiki Page](https://github.com/neelp03/ThrottleX/wiki/Installation-and-Setup)** for complete instructions.

---

## **Usage**

### Example Initialization

```go
import (
    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)

func main() {
    redisStore := store.NewRedisStore(redisClient)
    limiter, err := ratelimiter.NewLeakyBucketLimiter(redisStore, 100, time.Second, 5) // Capacity: 100, Leak rate: 1 req/sec, Concurrency: 5
    if err != nil {
        log.Fatalf("Failed to initialize rate limiter: %v", err)
    }

    if allowed, _ := limiter.Allow("client-id"); allowed {
        // Process request
    }
}
```

For more example integrations, visit the **[Examples Wiki Page](https://github.com/neelp03/ThrottleX/wiki/ThrottleX-Examples)**.

---

## **Rate Limiting Algorithms**

ThrottleX offers the following algorithms to adapt to various rate-limiting needs:

- **Fixed Window Limiter**: Limits requests within set time frames.
- **Sliding Window Limiter**: Smoothes out request patterns over sliding intervals.
- **Token Bucket Limiter**: Allows bursts while limiting sustained traffic.
- **Leaky Bucket Limiter** (new): Controls request processing rate with a concurrency limiter to prevent overloads.

For detailed information, see the **[Rate Limiting Algorithms Wiki Page](https://github.com/neelp03/ThrottleX/wiki/Rate-Limiting-Algorithms-in-ThrottleX)**.

---

## **Changelog**

Changes between releases are documented in the `CHANGELOG.md` file. View the **[Changelog](https://github.com/neelp03/ThrottleX/blob/main/CHANGELOG.md)** for a detailed list of updates and bug fixes.

---

## **Contributing**

To contribute:

1. **Fork the Repository**.
2. **Clone Your Fork**:

   ```bash
   git clone https://github.com/your-username/throttlex.git
   cd throttlex
   ```

3. **Create a Branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make Changes and Run Tests**:

   ```bash
   go test -race -v ./...
   ```

5. **Commit and Push**:

   ```bash
   git add .
   git commit -m "Add your feature"
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request** on the `main` branch.

---

## **License**

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

---

## **Acknowledgments**

ThrottleX was created to address the need for flexible, high-performance rate limiting in Go applications. Special thanks to the Go community for their guidance and contributions.

---

## **Contact**

For questions or support, please open an issue on the **[GitHub repository](https://github.com/neelp03/throttlex/issues)**.
