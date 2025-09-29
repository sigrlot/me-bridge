package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// BeutConsoleWriter implements zerolog.ConsoleWriter interface for Ethereum-style logging.
// It formats log entries to match the output style of Ethereum clients like Geth.
type BeutConsoleWriter struct {
	Out        io.Writer // Out specifies the output destination
	NoColor    bool      // NoColor disables colored output
	TimeFormat string    // TimeFormat specifies the time format (defaults to Ethereum style)
}

// Write formats and writes a log entry in Ethereum style.
// It parses JSON log entries and converts them to the format:
// LEVEL[timestamp] [component] message field=value
func (w BeutConsoleWriter) Write(p []byte) (n int, err error) {
	var evt map[string]any

	// Parse the JSON log entry
	if err := json.Unmarshal(p, &evt); err != nil {
		return w.Out.Write(p)
	}

	// Extract fields
	timestamp := getStringField(evt, "time", time.Now().Format("01-02|15:04:05.000"))
	level := getStringField(evt, "level", "INFO")
	message := getStringField(evt, "message", "")
	component := getStringField(evt, "component", "")

	// Format timestamp to Ethereum style (MM-dd|HH:mm:ss.SSS)
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		timestamp = t.Format("01-02|15:04:05.000")
	}

	// Apply color formatting if enabled
	colorStart, colorEnd := getColorCodes(level, w.NoColor)

	// Convert level to Ethereum style
	levelStr := convertLogLevel(level)

	// Build the log line in Ethereum style
	logLine := formatLogLine(levelStr, timestamp, component, message, colorStart, colorEnd)

	// Add extra fields (excluding standard ones)
	extraFields := extractExtraFields(evt)
	if len(extraFields) > 0 {
		logLine += " " + strings.Join(extraFields, " ")
	}

	logLine += "\n"

	return w.Out.Write([]byte(logLine))
}

// Helper function to get string field from parsed event
func getStringField(evt map[string]any, key, defaultValue string) string {
	if val, ok := evt[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// convertLogLevel returns the original log level in uppercase, truncated to 4 characters
func convertLogLevel(level string) string {
	levelStr := strings.ToUpper(level)
	if len(levelStr) > 4 {
		return levelStr[:4]
	}
	return levelStr
}

// getColorCodes returns ANSI color codes for the given log level
func getColorCodes(level string, noColor bool) (start, end string) {
	if noColor {
		return "", ""
	}

	switch level {
	case "TRAC": // TRACE truncated to 4 chars
		return "\033[36m", "\033[0m" // Cyan
	case "DEBUG": // DEBUG truncated to 4 chars
		return "\033[37m", "\033[0m" // White
	case "INFO":
		return "\033[32m", "\033[0m" // Green
	case "WARN", "WARNING": // WARN/WARNING truncated to 4 chars
		return "\033[33m", "\033[0m" // Yellow
	case "ERROR": // ERROR truncated to 4 chars
		return "\033[31m", "\033[0m" // Red
	case "FATAL", "PANI": // FATAL/PANIC truncated to 4 chars
		return "\033[35m", "\033[0m" // Magenta
	default:
		return "", ""
	}
}

// formatLogLine builds the main log line in Ethereum format
func formatLogLine(level, timestamp, component, message, colorStart, colorEnd string) string {
	if component != "" {
		// Format: LEVEL[timestamp] [component] message
		return fmt.Sprintf("%s%s[%s] [%s] %s%s",
			colorStart, level, timestamp, component, message, colorEnd)
	}
	// Format: LEVEL[timestamp] message
	return fmt.Sprintf("%s%s[%s] %s%s",
		colorStart, level, timestamp, message, colorEnd)
}

// extractExtraFields extracts additional fields from the log event
func extractExtraFields(evt map[string]any) []string {
	standardFields := map[string]bool{
		"time":      true,
		"level":     true,
		"message":   true,
		"component": true,
		"caller":    true,
	}

	extraFields := make([]string, 0)
	for key, value := range evt {
		if !standardFields[key] {
			extraFields = append(extraFields, fmt.Sprintf("%s=%v", key, value))
		}
	}
	return extraFields
}
