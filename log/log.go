package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zerolog.Logger to provide additional functionality
type Logger struct {
	config *Config
	logger zerolog.Logger
}

// NewLogger creates a new logger with given configuration
func NewLogger(config *Config) (*Logger, error) {
	// Set log level
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", config.Level, err)
	}
	zerolog.SetGlobalLevel(level)

	// Configure output writer
	var writer io.Writer
	switch strings.ToLower(config.Output) {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	case "file":
		if config.Path == "" {
			return nil, fmt.Errorf("path is required when output is file")
		}

		// Create directory if not exists
		dir := filepath.Dir(config.Path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		writer = &lumberjack.Logger{
			Filename:   config.Path,
			MaxSize:    1024,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		}
	default:
		return nil, fmt.Errorf("unsupported output type: %s", config.Output)
	}

	// Configure format
	var logger zerolog.Logger
	switch strings.ToLower(config.Format) {
	case "json":
		logger = zerolog.New(writer)
	case "console":
		consoleWriter := zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
			NoColor:    config.Output == "file", // Disable color for file output
		}
		logger = zerolog.New(consoleWriter)
	case "ethereum":
		// Use custom Ethereum-style formatter
		// First create JSON logger, then wrap with custom writer
		ethWriter := BeutConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
			NoColor:    config.Output == "file", // Disable color for file output

		}
		logger = zerolog.New(ethWriter)
	default:
		return nil, fmt.Errorf("unsupported format type: %s (supported: json, console, ethereum)", config.Format)
	}

	// Add timestamp
	logger = logger.With().Timestamp().Logger()

	// Add caller info if enabled
	if true {
		logger = logger.With().Caller().Logger()
	}

	return &Logger{
		config: config,
		logger: logger,
	}, nil
}

// WithComponent creates a new logger with component field
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{logger: l.logger.With().Str("component", component).Logger()}
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields ...map[string]any) *Logger {
	event := l.logger.With()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}

	return &Logger{logger: event.Logger()}
}

// WithError creates a new logger with error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{logger: l.logger.With().Err(err).Logger()}
}

// addFields adds fields to the raw logger
func (l *Logger) addFields(fields ...map[string]any) *zerolog.Logger {
	event := l.logger.With()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	logger := event.Logger()
	return &logger
}

// Trace logs a trace level message
func (l *Logger) Trace(msg string, fields ...map[string]any) {
	l.addFields(fields...).Trace().Msg(msg)
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string, fields ...map[string]any) {
	l.addFields(fields...).Debug().Msg(msg)
}

// Info logs an info level message
func (l *Logger) Info(msg string, fields ...map[string]any) {
	l.addFields(fields...).Info().Msg(msg)
}

// Warn logs a warn level message
func (l *Logger) Warn(msg string, fields ...map[string]any) {
	l.addFields(fields...).Warn().Msg(msg)
}

// Error logs an error level message
func (l *Logger) Error(msg string, fields ...map[string]any) {
	l.addFields(fields...).Error().Msg(msg)
}

// Fatal logs a fatal level message and exits
func (l *Logger) Fatal(msg string, fields ...map[string]any) {
	l.addFields(fields...).Fatal().Msg(msg)
}

// Panic logs a panic level message and panics
func (l *Logger) Panic(msg string, fields ...map[string]any) {
	l.addFields(fields...).Panic().Msg(msg)
}
