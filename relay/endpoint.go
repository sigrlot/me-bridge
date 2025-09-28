package relay

// InEndpoint 代表跨链桥的跨入端
type InEndpoint interface {
	InWatcher
	OutProcessor
}

// Watcher 订阅跨链消息
type InWatcher interface {
	// SubscribeToInMsgs 订阅入向消息；端点使用自己的配置
	// 来决定订阅什么。返回入向消息的通道。
	SubscribeToInMsgs() (<-chan InMsg, error)
}

// Processor 处理跨链消息
type OutProcessor interface {
	ProcessOutMsgs(msgs <-chan OutMsg) error
}

type OutEndpoint interface {
	InProcessor
	OutWatcher
}

// Processor 处理跨链消息
type InProcessor interface {
	ProcessInMsgs(msgs <-chan InMsg) error
}

// Watcher 订阅跨链消息
type OutWatcher interface {
	SubscribeToOutMsgs() (<-chan OutMsg, error)
	ConfirmOutMsgs(msgs []OutMsg) error
	SubscribeToBatchMsgs() (<-chan *BatchMsg, error)
}
