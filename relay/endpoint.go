package relay

import "context"

// InEndpoint 代表跨链桥的跨入端
type InEndpoint interface {
	InWatcher
	OutProcessor
}

// Watcher 订阅跨链消息
type InWatcher interface {
	FilterInMsgs(fromHeight, toHeight uint64) ([]InMsg, error)
	// SubscribeToInMsgs 订阅跨入消息，消息通过通道发送
	SubscribeToInMsgs(msgs chan InMsg) error
}

// Processor 处理跨链消息
type OutProcessor interface {
	FeeCalculator
	ProcessOutMsgs(msgs <-chan OutMsg) error
}

type OutEndpoint interface {
	InProcessor
	OutWatcher

	// 状态同步方法
	GetSequence() (uint64, uint64)  // 返回 (sequence, height)
	GetNonce(address string) uint64 // 获取 nonce
}

// Processor 处理跨链消息
type InProcessor interface {
	// ProcessInMsgs 处理跨链消息
	ProcessInMsgs(msgs <-chan InMsg) error

	// HandleError 处理端点级别的错误
	HandleError(ctx context.Context, err error, metadata map[string]interface{}) error
}

// Watcher 订阅跨链消息
type OutWatcher interface {
	SubscribeToOutMsgs(msgs chan OutMsg) error
	ConfirmOutMsgs(msgs []OutMsg) error
	SubscribeToBatchMsgs() (<-chan *BatchMsg, error)
}
