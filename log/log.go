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

// Logger 包装 zerolog.Logger 以提供额外功能
type Logger struct {
	config *Config
	logger zerolog.Logger
}

// NewLogger 使用给定配置创建新的日志器
func NewLogger(config *Config) (*Logger, error) {
	// 设置日志级别
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", config.Level, err)
	}
	zerolog.SetGlobalLevel(level)

	// 配置输出写入器
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

		// 如果目录不存在则创建
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

	// 配置格式
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
		// 使用自定义以太坊风格格式化器
		// 首先创建 JSON 日志器，然后用自定义写入器包装
		ethWriter := BeutConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
			NoColor:    config.Output == "file", // Disable color for file output

		}
		logger = zerolog.New(ethWriter)
	default:
		return nil, fmt.Errorf("unsupported format type: %s (supported: json, console, ethereum)", config.Format)
	}

	// 添加时间戳
	logger = logger.With().Timestamp().Logger()

	// 如果启用则添加调用者信息
	if true {
		logger = logger.With().Caller().Logger()
	}

	return &Logger{
		config: config,
		logger: logger,
	}, nil
}

// WithComponent 创建带有组件字段的新日志器
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{logger: l.logger.With().Str("component", component).Logger()}
}

// WithFields 创建带有额外字段的新日志器
func (l *Logger) WithFields(fields ...map[string]any) *Logger {
	event := l.logger.With()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}

	return &Logger{logger: event.Logger()}
}

// WithError 创建带有错误字段的新日志器
func (l *Logger) WithError(err error) *Logger {
	return &Logger{logger: l.logger.With().Err(err).Logger()}
}

// addFields 向原始日志器添加字段
func (l *Logger) addFields(event *zerolog.Event, fields ...map[string]any) *zerolog.Event {
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	return event
}

// Trace 记录跟踪级别消息
func (l *Logger) Trace(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Trace(), fields...).Msg(msg)
}

// Debug 记录调试级别消息
func (l *Logger) Debug(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Debug(), fields...).Msg(msg)
}

// Info 记录信息级别消息
func (l *Logger) Info(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Info(), fields...).Msg(msg)
}

// Warn 记录警告级别消息
func (l *Logger) Warn(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Warn(), fields...).Msg(msg)
}

// Error 记录错误级别消息
func (l *Logger) Error(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Error(), fields...).Msg(msg)
}

// Fatal 记录致命级别消息并退出
func (l *Logger) Fatal(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Fatal(), fields...).Msg(msg)
}

// Panic 记录恐慌级别消息并恐慌
func (l *Logger) Panic(msg string, fields ...map[string]any) {
	l.addFields(l.logger.Panic(), fields...).Msg(msg)
}
