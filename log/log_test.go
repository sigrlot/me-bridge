package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", config.Level)
	}
	if config.Format != "ethereum" {
		t.Errorf("Expected format 'ethereum', got '%s'", config.Format)
	}
	if config.Output != "stdout" {
		t.Errorf("Expected output 'stdout', got '%s'", config.Output)
	}
	if config.Path != "logs/app.log" {
		t.Errorf("Expected output 'stdout', got '%s'", config.Path)
	}
}

func TestNewLoggerWithInvalidLevel(t *testing.T) {
	config := &Config{
		Level:  "invalid",
		Format: "json",
		Output: "stdout",
	}

	_, err := NewLogger(config)
	if err == nil {
		t.Error("Expected error for invalid log level")
	}
}

func TestNewLoggerWithInvalidOutput(t *testing.T) {
	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "invalid",
	}

	_, err := NewLogger(config)
	if err == nil {
		t.Error("Expected error for invalid output type")
	}
}

func TestNewLoggerWithInvalidFormat(t *testing.T) {
	config := &Config{
		Level:  "info",
		Format: "invalid",
		Output: "stdout",
	}

	_, err := NewLogger(config)
	if err == nil {
		t.Error("Expected error for invalid format type")
	}
}

func TestNewLoggerFileOutput(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "file",
		Path:   logFile,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.Info("test message")

	// Check if file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

func TestLoggerWithComponent(t *testing.T) {
	var buf bytes.Buffer

	zlog := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{logger: zlog}

	logger.WithComponent("api").Info("component message")

	output := buf.String()
	if !strings.Contains(output, "component") {
		t.Error("Expected log to contain 'component' field")
	}
	if !strings.Contains(output, "api") {
		t.Error("Expected log to contain 'api' value")
	}
}

func TestLoggerWithFields(t *testing.T) {
	var buf bytes.Buffer

	// Create a logger that writes to buffer for testing
	zlog := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{logger: zlog}

	fields := map[string]any{
		"user_id": 123,
		"action":  "login",
	}

	logger.WithFields(fields).Info("user action")

	output := buf.String()
	if !strings.Contains(output, "user_id") {
		t.Error("Expected log to contain 'user_id' field")
	}
	if !strings.Contains(output, "action") {
		t.Error("Expected log to contain 'action' field")
	}
}

func TestLoggerWithError(t *testing.T) {
	var buf bytes.Buffer

	zlog := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{logger: zlog}

	testErr := errors.New("test error")
	logger.WithError(testErr).Error("error occurred")

	output := buf.String()
	if !strings.Contains(output, "error") {
		t.Error("Expected log to contain 'error' field")
	}
	if !strings.Contains(output, "test error") {
		t.Error("Expected log to contain error message")
	}
}

func TestAllLogLevels(t *testing.T) {
	var buf bytes.Buffer

	// Save and restore global level
	oldLevel := zerolog.GlobalLevel()
	defer zerolog.SetGlobalLevel(oldLevel)
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	zlog := zerolog.New(&buf).Level(zerolog.TraceLevel).With().Timestamp().Logger()
	logger := &Logger{logger: zlog}

	logger.Trace("trace message")
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	expectedMessages := []string{
		"trace message",
		"debug message",
		"info message",
		"warn message",
		"error message",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected log to contain '%s', but got output: %s", msg, output)
		}
	}
}

func TestGlobalLogger(t *testing.T) {
	// Test that global logger is initialized
	logger := GetLogger()
	if logger == nil {
		t.Error("Global logger should not be nil")
	}

	// Test global functions
	var buf bytes.Buffer
	zlog := zerolog.New(&buf)
	RootLogger = &Logger{logger: zlog}

	Info("global info message")
	output := buf.String()

	if !strings.Contains(output, "global info message") {
		t.Error("Expected global info message in log output")
	}
}

func TestSetLogger(t *testing.T) {
	// Test init with nil config (should use default)
	err := SetLogger(nil)
	if err != nil {
		t.Errorf("Init with nil config should not fail: %v", err)
	}

	// Test init with custom config
	tempDir := t.TempDir()
	config := &Config{
		Level:  "debug",
		Format: "json",
		Output: "file",
		Path:   filepath.Join(tempDir, "init_test.log"),
	}

	err = SetLogger(config)
	if err != nil {
		t.Errorf("Init with custom config should not fail: %v", err)
	}
}

func TestJSONOutput(t *testing.T) {
	var buf bytes.Buffer

	// We need to manually create a JSON logger for testing
	zlog := zerolog.New(&buf).With().Timestamp().Logger()
	logger := &Logger{logger: zlog}

	logger.Info("json test message", map[string]any{"test_field": "test_value"})

	output := buf.String()

	// Verify it's valid JSON
	var jsonData map[string]any
	if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}

	// Check specific fields
	if jsonData["test_field"] != "test_value" {
		t.Error("Expected test_field to be 'test_value'")
	}
	if jsonData["message"] != "json test message" {
		t.Error("Expected message to be 'json test message'")
	}
}

func TestConcurrentLogging(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "file",
		Path:   filepath.Join(tempDir, "concurrent_test.log"),
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Run concurrent logging
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Info("concurrent message", map[string]any{"goroutine": id, "iteration": j})
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Just verify the file exists and is not empty
	info, err := os.Stat(config.Path)
	if err != nil {
		t.Errorf("Log file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Log file should not be empty")
	}
}

func TestLogRotation(t *testing.T) {
	tempDir := t.TempDir()
	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "file",
		Path:   filepath.Join(tempDir, "rotation_test.log"),
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Write some log messages
	for i := 0; i < 1000; i++ {
		logger.Info("This is a test message for log rotation functionality", map[string]any{"iteration": i})
	}

	// Check if log file exists
	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		t.Error("Log file should exist")
	}
}

func TestEthereumFormat(t *testing.T) {
	var buf bytes.Buffer

	// Create ethereum-style logger that writes to buffer
	ethWriter := BeutConsoleWriter{
		Out:     &buf,
		NoColor: true, // Disable color for testing
	}

	zlog := zerolog.New(ethWriter).With().Timestamp().Logger()
	logger := &Logger{logger: zlog}

	logger.WithComponent("p2p").Info("Started P2P networking")
	logger.WithComponent("blockchain").Info("Imported new block", map[string]any{"number": 12345})

	output := buf.String()

	// Check ethereum-style format
	if !strings.Contains(output, "INFO[") {
		t.Error("Expected ethereum-style INFO level format")
	}
	if !strings.Contains(output, "[p2p]") {
		t.Error("Expected p2p component in ethereum format")
	}
	if !strings.Contains(output, "[blockchain]") {
		t.Error("Expected blockchain component in ethereum format")
	}
	if !strings.Contains(output, "Started P2P networking") {
		t.Error("Expected message in ethereum format")
	}
	if !strings.Contains(output, "number=12345") {
		t.Error("Expected field in ethereum format")
	}
}

// Performance test to ensure we maintain good performance
func TestPerformance(t *testing.T) {
	var buf bytes.Buffer
	zlog := zerolog.New(&buf).Level(zerolog.InfoLevel)
	logger := &Logger{logger: zlog}

	start := time.Now()

	// Log 10000 messages
	for i := 0; i < 10000; i++ {
		logger.Info("performance test message", map[string]any{"iteration": i})
	}

	duration := time.Since(start)

	// This is a rough performance check - adjust based on your requirements
	if duration > time.Second {
		t.Errorf("Logging 10000 messages took too long: %v", duration)
	}
}
