package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *slog.Logger
)

// Init initializes the structured logger with JSON handler
func Init() {
	once.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger = slog.New(handler)
		slog.SetDefault(logger)
	})
}

// Get returns the singleton logger instance
func Get() *slog.Logger {
	if logger == nil {
		Init()
	}
	return logger
}

// Info logs info level message
func Info(msg string, args ...interface{}) {
	Get().Info(msg, args...)
}

// Debug logs debug level message
func Debug(msg string, args ...interface{}) {
	Get().Debug(msg, args...)
}

// Warn logs warning level message
func Warn(msg string, args ...interface{}) {
	Get().Warn(msg, args...)
}

// Error logs error level message
func Error(msg string, args ...interface{}) {
	Get().Error(msg, args...)
}

// WithRequest adds request context to logger
func WithRequest(requestID, method, path string, userID uint) *slog.Logger {
	return logger.With(
		slog.String("request_id", requestID),
		slog.String("method", method),
		slog.String("path", path),
		slog.Uint64("user_id", uint64(userID)),
	)
}
