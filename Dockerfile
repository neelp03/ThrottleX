# Build Stage
FROM golang:1.23-alpine AS build

# Set the working directory
WORKDIR /app

# Copy go mod and sum files for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o throttlex ./cmd/throttlex

# Final Stage
FROM alpine:latest

# Set working directory in the final image
WORKDIR /root/

# Copy the built Go binary from the build stage
COPY --from=build /app/throttlex .

# Expose port 8080 for the API and 2112 for Prometheus metrics
EXPOSE 8080
EXPOSE 2112

# Run the Go app
CMD ["./throttlex"]
