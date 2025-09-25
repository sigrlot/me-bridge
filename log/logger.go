package log

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Global logger instance
var RootLogger *Logger

// DefaultConfig returns default logging configuration
func DefaultConfig() *Config {
	return &Config{
		Level:  "info",
		Format: "ethereum",
		Output: "stdout",
		Path:   "logs/app.log",
	}
}

// Init initializes the global logger with given configuration
func SetLogger(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// Set error stack trace marshaler
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Create logger
	logger, err := NewLogger(config)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	RootLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if RootLogger == nil {
		// Initialize with default config if not initialized
		_ = SetLogger(nil)
	}
	return RootLogger
}

// Global convenience functions
func WithComponent(component string) *Logger {
	return GetLogger().WithComponent(component)
}

func WithFields(fields ...map[string]any) *Logger {
	return GetLogger().WithFields(fields...)
}

func WithError(err error) *Logger {
	return GetLogger().WithError(err)
}

func Trace(msg string, fields ...map[string]any) {
	GetLogger().Trace(msg, fields...)
}

func Debug(msg string, fields ...map[string]any) {
	GetLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...map[string]any) {
	GetLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...map[string]any) {
	GetLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...map[string]any) {
	GetLogger().Error(msg, fields...)
}

func Fatal(msg string, fields ...map[string]any) {
	GetLogger().Fatal(msg, fields...)
}

func Panic(msg string, fields ...map[string]any) {
	GetLogger().Panic(msg, fields...)
}
