package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

// LoadConfig loads the environment variables from .env or system environment
func LoadConfig() {
    // Load from .env file if present
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Check for required variables
    if os.Getenv("REDIS_HOST") == "" || os.Getenv("PORT") == "" {
        log.Fatal("Missing required environment variables: REDIS_HOST, PORT")
    }
}
