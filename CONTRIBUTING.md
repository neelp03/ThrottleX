# Contributing to Throttlex

Thank you for considering contributing to **Throttlex**! Your contributions are highly appreciated and help improve this project for everyone.

---

## Table of Contents

- [Getting Started](#getting-started)
  - [Fork the Repository](#fork-the-repository)
  - [Clone Your Fork](#clone-your-fork)
  - [Set Upstream Remote](#set-upstream-remote)
- [Setting Up the Development Environment](#setting-up-the-development-environment)
  - [Prerequisites](#prerequisites)
  - [Installing Dependencies](#installing-dependencies)
- [Code Contribution Guidelines](#code-contribution-guidelines)
  - [Coding Standards](#coding-standards)
  - [Formatting and Linting](#formatting-and-linting)
  - [Commit Messages](#commit-messages)
- [Running Tests](#running-tests)
  - [Unit Tests](#unit-tests)
  - [Integration Tests](#integration-tests)
- [Submitting Changes](#submitting-changes)
  - [Pull Request Guidelines](#pull-request-guidelines)
- [Issue Reporting](#issue-reporting)
- [Code of Conduct](#code-of-conduct)
- [License](#license)

---

## Getting Started

### Fork the Repository

Click the **Fork** button at the top right of the [repository page](https://github.com/neelp03/throttlex) to create your own copy of the repository.

### Clone Your Fork

Clone your forked repository to your local machine:

```bash
git clone https://github.com/yourusername/throttlex.git
cd throttlex
```

Replace `yourusername` with your GitHub username.

### Set Upstream Remote

Set the original repository as the upstream remote:

```bash
git remote add upstream https://github.com/neelp03/throttlex.git
```

This allows you to keep your fork up-to-date with the latest changes.

---

## Setting Up the Development Environment

### Prerequisites

- **Go**: Ensure you have Go version **1.21.x** or later installed. Download it from [golang.org](https://golang.org/dl/).

- **Git**: Install Git and configure it with your user name and email.

- **Redis**: If you plan to work on features involving Redis, install Redis locally or run it via Docker.

### Installing Dependencies

Navigate to the project directory and download the module dependencies:

```bash
go mod download
```

---

## Code Contribution Guidelines

### Coding Standards

- **Follow Go Conventions**: Adhere to the standard Go coding conventions. Refer to [Effective Go](https://golang.org/doc/effective_go.html) for guidance.

- **Comment Your Code**: Document all exported functions, types, and packages with clear comments.

- **Write Idiomatic Go**: Use idiomatic Go patterns and best practices.

### Formatting and Linting

- **Code Formatting**: Use `gofmt -s` to format your code.

  ```bash
  gofmt -s -w .
  ```

- **Linting**: Use `golangci-lint` to check for code issues.

  ```bash
  golangci-lint run
  ```

- **Editor Configuration**: Configure your editor or IDE to format code on save and highlight linting issues.

### Commit Messages

- **Use Clear Messages**: Write concise and descriptive commit messages.

- **Commit Structure**:

  ```
  Subject: Brief description of the change (max 50 characters)

  Optional detailed explanation, wrapping at 72 characters.
  ```

- **Example**:

  ```
  Add sliding window rate limiter

  Implemented a new sliding window rate limiter to provide smoother rate limiting over time windows.
  ```

---

## Running Tests

### Unit Tests

Run unit tests with:

```bash
go test -race -v ./...
```

### Integration Tests

Integration tests may require external services like Redis.

- **Start Redis**: Ensure Redis is running locally on `localhost:6379`. You can use Docker:

  ```bash
  docker run -p 6379:6379 -d redis:6.2
  ```

- **Run Integration Tests**:

  ```bash
  go test -tags=integration -race -v ./...
  ```

- **Environment Variables**: If Redis is on a different host or port, set `REDIS_ADDR`:

  ```bash
  REDIS_ADDR=your_redis_host:6379 go test -tags=integration -race -v ./...
  ```

---

## Submitting Changes

### Pull Request Guidelines

1. **Sync with Upstream**:

   ```bash
   git checkout main
   git fetch upstream
   git merge upstream/main
   ```

2. **Create a Branch**:

   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**: Implement your feature or bug fix, and ensure code is well-documented.

4. **Run Tests**: Verify that all tests pass.

   ```bash
   go test -race -v ./...
   ```

5. **Commit Changes**:

   ```bash
   git add .
   git commit -m "Describe your changes"
   ```

6. **Push to Your Fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request**:

   - Navigate to your fork on GitHub.
   - Click on **New pull request**.
   - Provide a clear title and description.
   - Reference any related issues (e.g., `Closes #123`).

---

## Issue Reporting

If you encounter a bug or have a feature request, please open an issue on the [GitHub repository](https://github.com/neelp03/throttlex/issues).

- **Search Existing Issues**: Before opening a new issue, please check if it has already been reported.

- **Provide Detailed Information**: Include steps to reproduce, expected vs. actual behavior, and any relevant logs or screenshots.

---

## Code of Conduct

We are committed to fostering an open and welcoming environment. Please read and adhere to our [Code of Conduct](CODE_OF_CONDUCT.md).

---

## License

By contributing to Throttlex, you agree that your contributions will be licensed under the [Apache 2.0 License](LICENSE).

---

Thank you for contributing to Throttlex!

---
