package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// BeutConsoleWriter 实现了 zerolog.ConsoleWriter 接口，用于以太坊风格的日志输出。
// 它将日志条目格式化为类似 Geth 等以太坊客户端的输出风格。
type BeutConsoleWriter struct {
	Out        io.Writer // Out 指定输出目标
	NoColor    bool      // NoColor 禁用彩色输出
	TimeFormat string    // TimeFormat 指定时间格式（默认为以太坊风格）
}

// Write 以以太坊风格格式化并写出一条日志。
// 它会解析 JSON 日志条目并转换为如下格式：
// LEVEL[时间戳] [组件] 消息 键=值
func (w BeutConsoleWriter) Write(p []byte) (n int, err error) {
	var evt map[string]any

	// 解析 JSON 日志条目
	if err := json.Unmarshal(p, &evt); err != nil {
		return w.Out.Write(p)
	}

	// 提取字段
	timestamp := getStringField(evt, "time", time.Now().Format("01-02|15:04:05.000"))
	level := getStringField(evt, "level", "INFO")
	message := getStringField(evt, "message", "")
	component := getStringField(evt, "component", "")

	// 将时间戳格式化为以太坊风格 (MM-dd|HH:mm:ss.SSS)
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		timestamp = t.Format("01-02|15:04:05.000")
	}

	// 如启用则应用颜色格式
	colorStart, colorEnd := getColorCodes(level, w.NoColor)

	// 将日志级别转换为以太坊风格
	levelStr := convertLogLevel(level)

	// 构造以太坊风格的日志行
	logLine := formatLogLine(levelStr, timestamp, component, message, colorStart, colorEnd)

	// 添加额外字段（排除标准字段）
	extraFields := extractExtraFields(evt)
	if len(extraFields) > 0 {
		logLine += " " + strings.Join(extraFields, " ")
	}

	logLine += "\n"

	return w.Out.Write([]byte(logLine))
}

// 从解析后的事件中获取字符串字段的辅助函数
func getStringField(evt map[string]any, key, defaultValue string) string {
	if val, ok := evt[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// convertLogLevel 将原始日志级别转为大写，并截断为 4 个字符
func convertLogLevel(level string) string {
	levelStr := strings.ToUpper(level)
	if len(levelStr) > 4 {
		return levelStr[:4]
	}
	return levelStr
}

// getColorCodes 根据日志级别返回 ANSI 颜色代码
func getColorCodes(level string, noColor bool) (start, end string) {
	if noColor {
		return "", ""
	}

	switch level {
	case "TRAC": // TRACE 截断为 4 个字符
		return "\033[36m", "\033[0m" // 青色
	case "DEBUG": // DEBUG 截断为 4 个字符
		return "\033[37m", "\033[0m" // 白色
	case "INFO":
		return "\033[32m", "\033[0m" // 绿色
	case "WARN", "WARNING": // WARN/WARNING 截断为 4 个字符
		return "\033[33m", "\033[0m" // 黄色
	case "ERROR": // ERROR 截断为 4 个字符
		return "\033[31m", "\033[0m" // 红色
	case "FATAL", "PANI": // FATAL/PANIC 截断为 4 个字符
		return "\033[35m", "\033[0m" // 品红色
	default:
		return "", ""
	}
}

// formatLogLine 以以太坊格式构建主日志行
func formatLogLine(level, timestamp, component, message, colorStart, colorEnd string) string {
	if component != "" {
		// 格式：LEVEL[时间戳] [组件] 消息
		return fmt.Sprintf("%s%s[%s] [%s] %s%s",
			colorStart, level, timestamp, component, message, colorEnd)
	}
	// 格式：LEVEL[时间戳] 消息
	return fmt.Sprintf("%s%s[%s] %s%s",
		colorStart, level, timestamp, message, colorEnd)
}

// extractExtraFields 从日志事件中提取额外字段
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
