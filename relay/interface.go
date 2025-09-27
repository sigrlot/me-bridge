package relay

// Endpoint 代表跨链桥的一端
type Endpoint interface {
	Work(inChan chan *Message, outChan chan *Message) error
}

// Watcher 订阅跨链消息
type Watcher interface {
	SubscribeToRelayMsgs(address string) (<-chan Message, error)
}

// Processor 处理跨链消息
type Processor interface {
	ProcessRelayMsgs(msgs <-chan Message) error
}

// Checker 检查跨链桥状态，并发出预警信息
type Checker interface {
	BalanceWarning() error
}

type Message interface {
	GetTxHash() string
	GetSender() string
	GetReceiver() string
	GetAmount() string
}
