package types

import "context"

type ErrorHandler interface {
	HandleError(ctx context.Context, err error, metadata map[string]interface{}) error
}
