package log

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Global logger instance
var globalLogger *Logger

// DefaultConfig returns default logging configuration
func DefaultConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:  "info",
		Format: "ethereum",
		Output: "stdout",
		Path:   "logs/app.log",
	}
}

// Init initializes the global logger with given configuration
func SetLogger(config *LoggerConfig) error {
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

	globalLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if globalLogger == nil {
		// Initialize with default config if not initialized
		_ = SetLogger(nil)
	}
	return globalLogger
}

// Global convenience functions
func Trace(msg string) {
	GetLogger().Trace(msg)
}

func Tracef(format string, args ...interface{}) {
	GetLogger().Tracef(format, args...)
}

func Debug(msg string) {
	GetLogger().Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(msg string) {
	GetLogger().Info(msg)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(msg string) {
	GetLogger().Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(msg string) {
	GetLogger().Error(msg)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(msg string) {
	GetLogger().Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func Panic(msg string) {
	GetLogger().Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	GetLogger().Panicf(format, args...)
}

// WithFields creates a new logger with additional fields
func WithFields(fields map[string]interface{}) *Logger {
	return GetLogger().WithFields(fields)
}

// WithField creates a new logger with additional field
func WithField(key string, value interface{}) *Logger {
	return GetLogger().WithField(key, value)
}

// WithComponent creates a new logger with component field
func WithComponent(component string) *Logger {
	return GetLogger().WithComponent(component)
}

// WithError creates a new logger with error field
func WithError(err error) *Logger {
	return GetLogger().WithError(err)
}
