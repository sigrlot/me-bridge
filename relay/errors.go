package relay

import "errors"

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
