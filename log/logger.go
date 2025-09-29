package log

import (
	"fmt"
)

// 全局日志器实例
var RootLogger *Logger

func init() {
	// 初始化全局日志器
	logger, err := NewLogger("info", "beut", "stdout", "logs/app.log")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize global logger: %v", err))
	}
	RootLogger = logger
}

// SetRootLogger 使用给定配置初始化全局日志器
func SetRootLogger(logger *Logger) {
	RootLogger = logger
}

// GetRootLogger 返回全局日志器实例
func GetRootLogger() *Logger {
	return RootLogger
}

// 全局便利函数
func WithComponent(component string) *Logger {
	return RootLogger.WithComponent(component)
}

func WithFields(fields ...map[string]any) *Logger {
	return RootLogger.WithFields(fields...)
}

func WithError(err error) *Logger {
	return RootLogger.WithError(err)
}

func Trace(msg string, fields ...map[string]any) {
	RootLogger.Trace(msg, fields...)
}

func Debug(msg string, fields ...map[string]any) {
	RootLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...map[string]any) {
	RootLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...map[string]any) {
	RootLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...map[string]any) {
	RootLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...map[string]any) {
	RootLogger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...map[string]any) {
	RootLogger.Panic(msg, fields...)
}
