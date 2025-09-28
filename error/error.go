package error

import "context"

type ErrorHandler interface {
	HandleError(ctx context.Context, err error, metadata map[string]interface{}) error
}

// 错误处理策略
// type ErrorStrategy int
