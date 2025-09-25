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
	config *LoggerConfig
	logger zerolog.Logger
}

// NewLogger creates a new logger with given configuration
func NewLogger(config *LoggerConfig) (*Logger, error) {
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
func (l *Logger) withFields(fields map[string]interface{}) *zerolog.Logger {
	event := l.logger.With()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	logger := event.Logger()
	return &logger
}

// WithError creates a new logger with error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{logger: l.logger.With().Err(err).Logger()}
}

// Trace logs a trace level message
func (l *Logger) Trace(msg string, fields map[string]interface{}) {
	l.withFields(fields).Trace().Msg(msg)
}

// Tracef logs a trace level message with formatting
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logger.Trace().Msgf(format, args...)
}

// Debug logs a debug level message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a debug level message with formatting
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Info logs an info level message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs an info level message with formatting
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warn logs a warn level message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a warn level message with formatting
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Error logs an error level message
func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

// Errorf logs an error level message with formatting
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Error().Msgf(format, args...)
}

// Fatal logs a fatal level message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a fatal level message with formatting and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// Panic logs a panic level message and panics
func (l *Logger) Panic(msg string) {
	l.logger.Panic().Msg(msg)
}

// Panicf logs a panic level message with formatting and panics
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panic().Msgf(format, args...)
}
