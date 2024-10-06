package database

import (
    "github.com/go-redis/redis/v8"
    "log"
    "os"
    "context"
)

// RedisClient holds the Redis client connection instance.
// It is initialized by the `InitializeRedis` function and can be used throughout the application.
var RedisClient *redis.Client

// InitializeRedis establishes a connection to the Redis server.
//
// This function configures and connects the application to a Redis instance using
// the connection details specified in environment variables. The connection is
// established via the `redis/v8` package.
//
// It uses the following environment variables:
//   - REDIS_HOST: The address of the Redis server (required).
//   - REDIS_PASSWORD: The password for Redis (optional, can be empty if no password is set).
//
// After setting up the connection, the function sends a "PING" to verify that the
// connection to Redis is active. If the connection fails, the application will log
// a fatal error and terminate.
//
// Logs:
//   - "Connected to Redis: PONG" if the connection is successful.
//   - "Unable to connect to Redis" followed by the error message if the connection fails.
//
// Example usage:
//   database.InitializeRedis() // Initializes the Redis connection
//
// Fatal Error:
//   The application will terminate if the Redis connection fails.
func InitializeRedis() {
    // Create a new Redis client using options from environment variables
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST"),      // Redis server address
        Password: os.Getenv("REDIS_PASSWORD"),  // Redis password (optional)
        DB:       0,                            // Use default Redis DB
    })

    // Test the connection by sending a PING command to Redis
    pong, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatal("Unable to connect to Redis: ", err)
    }

    // Log successful connection
    log.Println("Connected to Redis: ", pong)
}
