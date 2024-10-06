package utils

import (
    "github.com/sirupsen/logrus"
    "os"
)

// log is the global logger instance initialized with Logrus.
var log = logrus.New()

// InitializeLogger sets up the logging level and format based on environment variables.
//
// This function configures the log level for the application using the `LOG_LEVEL`
// environment variable. It supports the following log levels:
//   - "debug": Logs debug-level and higher messages.
//   - "info": Logs info-level and higher messages (default).
// If the environment variable is not set, the default log level is "info".
//
// The log format is set to JSON using Logrus's `JSONFormatter`, which is useful
// for structured logging in production environments.
//
// Environment Variables:
//   - LOG_LEVEL: The desired log level ("debug", "info").
//
// Example usage:
//   utils.InitializeLogger() // Initializes the logger with the specified level and format.
func InitializeLogger() {
    // Get the log level from environment variables
    logLevel := os.Getenv("LOG_LEVEL")

    // Set the log level based on the environment variable
    switch logLevel {
    case "debug":
        log.SetLevel(logrus.DebugLevel)
    case "info":
        log.SetLevel(logrus.InfoLevel)
    default:
        log.SetLevel(logrus.InfoLevel)
    }

    // Set the log format to JSON for structured logging
    log.SetFormatter(&logrus.JSONFormatter{})
}

// LogError logs error messages along with the associated error object.
//
// This function logs error messages using the `logrus.Error` method, which includes
// both a custom error message and an associated error object for additional context.
// The log entry includes the error details in structured JSON format.
//
// Params:
//   - message: string - A descriptive error message to be logged.
//   - err: error - The error object that provides further context about the failure.
//
// Example usage:
//   utils.LogError("Failed to connect to Redis", err)
func LogError(message string, err error) {
    log.WithFields(logrus.Fields{
        "error": err,
    }).Error(message)
}

// LogInfo logs informational messages.
//
// This function logs informational messages using the `logrus.Info` method.
// It is useful for logging normal application events or status updates.
//
// Params:
//   - message: string - The informational message to be logged.
//
// Example usage:
//   utils.LogInfo("ThrottleX service has started successfully")
func LogInfo(message string) {
    log.Info(message)
}
