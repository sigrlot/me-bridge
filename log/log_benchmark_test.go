package log

import (
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// BenchmarkZerologJSON benchmarks JSON logging
func BenchmarkZerologJson(b *testing.B) {
	// Use direct zerolog to avoid file rotation issues in benchmarks
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("benchmark test message")
		}
	})
}

// BenchmarkZerologWithFields benchmarks logging with fields
func BenchmarkZerologJsonWithFields(b *testing.B) {
	// Use direct zerolog to avoid file rotation issues in benchmarks
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.WithFields(map[string]interface{}{
				"component":  "log",
				"user_id":    12345,
				"request_id": "req-abc-123",
				"timestamp":  time.Now(),
			}).Info("benchmark test message")
		}
	})
}

// BenchmarkZerologConsole benchmarks console logging
func BenchmarkZerologConsole(b *testing.B) {
	// Use console writer with io.Discard to avoid file rotation issues
	consoleWriter := zerolog.ConsoleWriter{
		Out:        io.Discard,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("benchmark test message")
		}
	})
}

// BenchmarkZerologConsole benchmarks console logging with fields
func BenchmarkZerologConsoleWithFields(b *testing.B) {
	// Use console writer with io.Discard to avoid file rotation issues
	consoleWriter := zerolog.ConsoleWriter{
		Out:        io.Discard,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.WithFields(map[string]interface{}{
				"component":  "log",
				"user_id":    12345,
				"request_id": "req-abc-123",
				"timestamp":  time.Now(),
			}).Info("benchmark test message")
		}
	})
}

// BenchmarkZerologBeut benchmarks Ethereum-style logging
func BenchmarkZerologBeut(b *testing.B) {
	// Use Ethereum-style writer with io.Discard to avoid file rotation issues
	ethWriter := BeutConsoleWriter{
		Out:        io.Discard,
		NoColor:    true,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(ethWriter).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("benchmark test message")
		}
	})
}

// BenchmarkZerologBeutWithModule benchmarks Ethereum-style logging with module
func BenchmarkZerologBeutWithFields(b *testing.B) {
	// Use Ethereum-style writer with io.Discard to avoid file rotation issues
	ethWriter := BeutConsoleWriter{
		Out:        io.Discard,
		NoColor:    true,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(ethWriter).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.WithFields(map[string]interface{}{
				"component":  "log",
				"user_id":    12345,
				"request_id": "req-abc-123",
				"timestamp":  time.Now(),
			}).Info("benchmark test message")
		}
	})
}
