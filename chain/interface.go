package chain

// Quirer 定义了中继器请求接口
type Quirer interface {
	GetNonce() uint64
	IncNonce() uint64
	RecordTx(hash string, nonce uint64) error
}

// Sender 定义了中继器交易发送接口
type Sender interface {
	Send() error
}

// Subscriber 定义了中继器订阅跨链消息的接口
type Subscriber interface {
	Subscribe() error
}

// Contracter 定义了中继器合约接口
type Contracter interface {
	Deploy() error
}

// Moduler 定义了中继器模块接口
type Moduler interface{}
