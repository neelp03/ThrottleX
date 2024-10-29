# Changelog

All notable changes to this project will be documented in this file.

## [v1.0.0-rc2] - 2024-10-29
### Added
- **Leaky Bucket Algorithm**: Introduced the Leaky Bucket rate-limiting algorithm to handle high-throughput scenarios, allowing requests to leak at a fixed rate.
- **Concurrency Control**: Added a concurrency limiter for the Leaky Bucket, preventing overloads by limiting the number of concurrent requests.
- **Unit Tests**: Expanded the testing suite to cover the Leaky Bucket and Concurrency Limiter, focusing on high concurrency and edge cases.
- **Version Retractions**: Added retractions for versions `v1.0.0` and `v1.0.1`, marking them as unavailable for improved backward compatibility in future releases.

### Improved
- **Redis Optimization (Preliminary)**: Prepared groundwork for Redis pipelining, setting the stage for more efficient backend operations and latency reduction.
- **Documentation**: Updated README and code comments to reflect the new Leaky Bucket algorithm and added examples for better clarity.

### Fixed
- **Concurrency Bug**: Resolved an issue where the semaphore for concurrency control could underflow, ensuring smoother performance during heavy loads.
- **Cleanup Routine**: Refined mutex cleanup process to optimize memory usage during high-traffic conditions, improving overall system stability.

## [v1.0.0-rc1] - 2024-10-16
### Initial Pre-release
- Introduced core functionality for ThrottleX, including Fixed Window, Sliding Window, and Token Bucket rate-limiting algorithms.
- Basic concurrency control and Redis integration for distributed environments.
- Initial testing suite for core rate-limiting algorithms.

---

