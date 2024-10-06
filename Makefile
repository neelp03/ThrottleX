build:
    go build -o throttlex cmd/throttlex/main.go

test:
    go test ./... -v

run:
    go run cmd/throttlex/main.go
