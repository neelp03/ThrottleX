
![image-1000x420 (7)](https://github.com/user-attachments/assets/55970be3-9e23-4613-b7ca-d58f9d73e0ed)

---

[![CI](https://github.com/neelp03/throttlex/actions/workflows/ci.yml/badge.svg)](https://github.com/neelp03/throttlex/actions/workflows/ci.yml)
[![Coverage Status](https://codecov.io/gh/neelp03/throttlex/branch/main/graph/badge.svg)](https://codecov.io/gh/neelp03/throttlex)
[![Go Report Card](https://goreportcard.com/badge/github.com/neelp03/throttlex?v=1)](https://goreportcard.com/report/github.com/neelp03/throttlex)
[![GoDoc](https://godoc.org/github.com/neelp03/throttlex?status.svg)](https://godoc.org/github.com/neelp03/throttlex)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

--- 

## **ThrottleX: Scalable Rate Limiting for Go APIs**

Welcome to ThrottleX, a powerful and flexible rate-limiting library for Go! 🚀
ThrottleX is designed to provide multiple rate limiting algorithms, easy integration, and scalable storage backends for your APIs.

For complete documentation, examples, and detailed setup instructions, please visit the **[ThrottleX Wiki](https://github.com/neelp03/ThrottleX/wiki)**.

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

- **Highly Configurable**:
  - Customize limits, window sizes, and keys.
  - Support for dynamic configurations.

- **Future Expansion**:
  - Upcoming support for additional rate limiting policies like Leaky Bucket, Concurrency Limit, and more.

---

## **Installation**

To install Throttlex, use `go get`:

```bash
go get -u github.com/neelp03/throttlex
```

For detailed setup instructions, refer to the **[Installation and Setup Wiki Page](https://github.com/neelp03/ThrottleX/wiki/Installation-and-Setup)**.

---

## **Usage**

Import the package into your Go project:

```go
import (
    "github.com/neelp03/throttlex/ratelimiter"
    "github.com/neelp03/throttlex/store"
)
```

For full examples of integrating ThrottleX with REST, gRPC, and GraphQL APIs, please refer to the **[Examples Wiki Page](https://github.com/neelp03/ThrottleX/wiki/ThrottleX-Examples)**.

---

## **Rate Limiting Algorithms**

ThrottleX currently supports the following rate limiting algorithms:

- **Fixed Window Limiter**
- **Sliding Window Limiter**
- **Token Bucket Limiter**

To learn more about these algorithms and how they work, visit the **[Rate Limiting Algorithms Wiki Page](https://github.com/neelp03/ThrottleX/wiki/Rate-Limiting-Algorithms-in-ThrottleX)**.

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

