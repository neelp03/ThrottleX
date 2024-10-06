# ThrottleX

**ThrottleX** is a high-performance, extensible Go library for rate limiting, designed to help you manage and control API traffic in your applications. It supports multiple rate-limiting algorithms and storage backends, making it suitable for a wide range of use cases, from single-server applications to distributed systems.

---

## Features

- **Multiple Rate Limiting Algorithms**:
  - **Fixed Window**
  - **Sliding Window**
  - **Token Bucket**

- **Pluggable Storage Backends**:
  - **Redis**: For distributed rate limiting across multiple instances.
  - **In-Memory**: For simplicity or testing purposes.

- **Easy Integration**:
  - Designed to integrate seamlessly with **REST APIs**, **gRPC services**, and more.
  - Provides clear interfaces and examples.

- **High Performance**:
  - Optimized algorithms for minimal latency.
  - Thread-safe implementations suitable for concurrent environments.

- **Extensible Design**:
  - Implement your own storage backends by adhering to the `Store` interface.
  - Customize rate-limiting parameters to fit your needs.

---

## Installation

Requires **Go 1.18** or newer.

```bash
go get github.com/neelp03/throttlex
```

---

## Quick Start

### **Import the Library**

```go
import (
    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)
```

### **Set Up a Rate Limiter**

```go
// Using an in-memory store
memStore := store.NewMemoryStore()

// Create a fixed window rate limiter
limiter := ratelimiter.NewFixedWindowLimiter(memStore, 100, time.Minute)
```

### **Use in Your Application**

```go
allowed, err := limiter.Allow("user-unique-key")
if err != nil {
    // Handle error
}

if !allowed {
    // Reject the request
}
```

---

## Examples

### **REST API Integration**

See [examples/rest_api/main.go](examples/rest_api/main.go) for a complete example using `net/http`.

### **gRPC Integration**

See [examples/grpc_api/main.go](examples/grpc_api/main.go) for how to use interceptors for rate limiting in gRPC.

---

## Documentation

- **[API Reference](https://pkg.go.dev/github.com/neelp03/throttlex)**: Full documentation of the library.
- **[Examples](examples/)**: Practical examples to help you get started.

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## License

This project is licensed under the **Apache 2.0** License - see the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or support, please open an issue on the GitHub repository.
