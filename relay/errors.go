package relay

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Common relay errors
var (
	ErrNonceNotFound       = errors.New("nonce not found")
	ErrInvalidMessage      = errors.New("invalid message")
	ErrChannelClosed       = errors.New("channel closed")
	ErrProcessingFailed    = errors.New("message processing failed")
	ErrSubscriptionFailed  = errors.New("subscription failed")
	ErrTransactionFailed   = errors.New("transaction failed")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrGasEstimationFailed = errors.New("gas estimation failed")
)

// ErrorAction 定义错误处理后的动作
type ErrorAction int

const (
	ActionRetry    ErrorAction = iota // 重试
	ActionIgnore                      // 忽略
	ActionEscalate                    // 升级到上一层处理
	ActionFatal                       // 致命错误，停止处理
)

// ErrorResult 错误处理结果
type ErrorResult struct {
	Action   ErrorAction            // 处理动作
	Delay    time.Duration          // 重试延迟（仅在 ActionRetry 时有效）
	Message  string                 // 处理消息
	Metadata map[string]interface{} // 附加信息
}

// ErrorLevel 定义错误级别
type ErrorLevel int

const (
	LevelClient   ErrorLevel = iota // 客户端级别
	LevelCluster                    // 集群级别
	LevelEndpoint                   // 端点级别
	LevelRelay                      // Relay级别
	LevelSigner                     // 签名器级别
)

// ErrorHandler 基础错误处理器
type ErrorHandler struct {
	MaxRetries int
	RetryDelay time.Duration
}

func NewErrorHandler(maxRetries int, retryDelay time.Duration) *ErrorHandler {
	return &ErrorHandler{
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}
}

// HandleError 基础错误处理逻辑
func (h *ErrorHandler) HandleError(ctx context.Context, err error, metadata map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// 根据错误类型决定处理策略
	action := h.classifyError(err)

	switch action {
	case ActionRetry:
		// 可以在当前层级重试
		return h.handleRetry(ctx, err, metadata)
	case ActionIgnore:
		// 忽略错误，继续处理
		return nil
	case ActionEscalate:
		// 升级到上一层处理
		return fmt.Errorf("escalating error from level %d: %w", h.Level, err)
	case ActionFatal:
		// 致命错误，停止处理
		return fmt.Errorf("fatal error at level %d: %w", h.Level, err)
	default:
		return err
	}
}

// classifyError 对错误进行分类
func (h *ErrorHandler) classifyError(err error) ErrorAction {
	errMsg := strings.ToLower(err.Error())

	// 可重试的错误
	if h.isRetryableError(errMsg) {
		return ActionRetry
	}

	// 可忽略的错误
	if h.isIgnorableError(errMsg) {
		return ActionIgnore
	}

	// 致命错误
	if h.isFatalError(errMsg) {
		return ActionFatal
	}

	// 默认升级处理
	return ActionEscalate
}

func (h *ErrorHandler) isRetryableError(errMsg string) bool {
	retryablePatterns := []string{
		"timeout", "network", "connection", "temporary",
		"nonce too low", "gas", "underpriced",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}
	return false
}

func (h *ErrorHandler) isIgnorableError(errMsg string) bool {
	ignorablePatterns := []string{
		"already known", "duplicate",
	}

	for _, pattern := range ignorablePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}
	return false
}

func (h *ErrorHandler) isFatalError(errMsg string) bool {
	fatalPatterns := []string{
		"invalid signature", "unauthorized", "permission denied",
	}

	for _, pattern := range fatalPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}
	return false
}

func (h *ErrorHandler) handleRetry(ctx context.Context, err error, metadata map[string]interface{}) error {
	retryCount := 0
	if count, ok := metadata["retryCount"].(int); ok {
		retryCount = count
	}

	if retryCount >= h.MaxRetries {
		return fmt.Errorf("max retries (%d) exceeded: %w", h.MaxRetries, err)
	}

	// 更新重试计数
	metadata["retryCount"] = retryCount + 1

	// 等待重试延迟
	select {
	case <-time.After(h.RetryDelay):
		return nil // 可以重试
	case <-ctx.Done():
		return ctx.Err()
	}
}
