package log

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// 全局日志器实例
var RootLogger *Logger

// DefaultConfig 返回默认日志配置
func DefaultConfig() *Config {
	return &Config{
		Level:  "info",
		Format: "ethereum",
		Output: "stdout",
		Path:   "logs/app.log",
	}
}

// SetLogger 使用给定配置初始化全局日志器
func SetLogger(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	// 设置错误堆栈跟踪序列化器
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// 创建日志器
	logger, err := NewLogger(config)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	RootLogger = logger
	return nil
}

// GetLogger 返回全局日志器实例
func GetLogger() *Logger {
	if RootLogger == nil {
		// 如果未初始化，则使用默认配置初始化
		_ = SetLogger(nil)
	}
	return RootLogger
}

// 全局便利函数
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
