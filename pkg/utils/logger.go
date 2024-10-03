package utils

import (
    "github.com/sirupsen/logrus"
    "os"
)

var log = logrus.New()

// InitializeLogger sets up the log level and format
func InitializeLogger() {
    logLevel := os.Getenv("LOG_LEVEL")
    switch logLevel {
    case "debug":
        log.SetLevel(logrus.DebugLevel)
    case "info":
        log.SetLevel(logrus.InfoLevel)
    default:
        log.SetLevel(logrus.InfoLevel)
    }

    log.SetFormatter(&logrus.JSONFormatter{})
}

// LogError logs error messages
func LogError(message string, err error) {
    log.WithFields(logrus.Fields{
        "error": err,
    }).Error(message)
}

// LogInfo logs informational messages
func LogInfo(message string) {
    log.Info(message)
}
