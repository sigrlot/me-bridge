package log

import (
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// BenchmarkZerologJSON 基准测试：JSON 日志
func BenchmarkZerologJson(b *testing.B) {
	// 直接使用 zerolog，避免基准测试中触发文件轮转
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("benchmark test message")
		}
	})
}

// BenchmarkZerologWithFields 基准测试：带字段的日志
func BenchmarkZerologJsonWithFields(b *testing.B) {
	// 直接使用 zerolog，避免基准测试中触发文件轮转
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	l := &Logger{logger: logger}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("benchmark test message", map[string]any{
				"component":  "log",
				"user_id":    12345,
				"request_id": "req-abc-123",
				"timestamp":  time.Now(),
			})
		}
	})
}

// BenchmarkZerologConsole 基准测试：控制台日志
func BenchmarkZerologConsole(b *testing.B) {
	// 使用 ConsoleWriter 且输出到 io.Discard，避免文件轮转影响
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

// BenchmarkZerologConsoleWithFields 基准测试：控制台日志（含字段）
func BenchmarkZerologConsoleWithFields(b *testing.B) {
	// 使用 ConsoleWriter 且输出到 io.Discard，避免文件轮转影响
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
			l.Info("benchmark test message", map[string]any{
				"component":  "log",
				"user_id":    12345,
				"request_id": "req-abc-123",
				"timestamp":  time.Now(),
			})
		}
	})
}

// BenchmarkZerologBeut 基准测试：以太坊风格日志
func BenchmarkZerologBeut(b *testing.B) {
	// 使用以太坊风格写入器并输出到 io.Discard，避免文件轮转影响
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

// BenchmarkZerologBeutWithFields 基准测试：以太坊风格日志（含字段）
func BenchmarkZerologBeutWithFields(b *testing.B) {
	// 使用以太坊风格写入器并输出到 io.Discard，避免文件轮转影响
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
			l.Info("benchmark test message", map[string]any{
				"component":  "log",
				"user_id":    12345,
				"request_id": "req-abc-123",
				"timestamp":  time.Now(),
			})
		}
	})
}
