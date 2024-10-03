package database

import (
    "github.com/go-redis/redis/v8"
    "log"
    "os"
    "context"
)

var RedisClient *redis.Client

// InitializeRedis sets up the Redis connection
func InitializeRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    // Test the connection
    pong, err := RedisClient.Ping(context.Background()).Result()
    if err != nil {
        log.Fatal("Unable to connect to Redis: ", err)
    }
    log.Println("Connected to Redis: ", pong)
}
