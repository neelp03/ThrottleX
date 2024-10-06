package config

import (
    "github.com/joho/godotenv"
    "log"
    "os"
)

// LoadConfig loads the configuration for ThrottleX from environment variables.
//
// This function attempts to load environment variables from a `.env` file if it exists.
// If the `.env` file is not found, the system environment variables will be used instead.
//
// The function ensures that the necessary environment variables `REDIS_HOST` and `PORT`
// are set, and it will log a fatal error and stop execution if any of the required variables
// are missing.
//
// Example:
//   config.LoadConfig() // Loads environment variables and performs checks
//
// Environment Variables:
//   - REDIS_HOST: Specifies the Redis server address.
//   - PORT: Defines the port on which the application will run.
//
// Logs:
//   - "No .env file found, using system environment variables" if no `.env` file is present.
//   - "Missing required environment variables: REDIS_HOST, PORT" if any required variables are absent.
//
// Fatal Error:
//   If the required environment variables are not set, the function will terminate the application.
func LoadConfig() {
    // Load from .env file if present
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    // Check for required environment variables
    if os.Getenv("REDIS_HOST") == "" || os.Getenv("PORT") == "" {
        log.Fatal("Missing required environment variables: REDIS_HOST, PORT")
    }
}
