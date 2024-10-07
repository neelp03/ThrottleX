.PHONY: test
test:
	# Use a tab instead of spaces here
	go test -race -v ./...

.PHONY: bench
bench:
	# Use a tab instead of spaces here
	go test -bench=. ./...

.PHONY: lint
lint:
	# Use a tab instead of spaces here
	golangci-lint run ./...

.PHONY: all
all: test lint bench

