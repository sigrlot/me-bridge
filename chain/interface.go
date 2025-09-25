package chain

// Client 封装了区块链client的差异性，提供了统一的跨链桥相关方法
type Client interface {
	// LatestHeight 获取最新区块高度
	LatestHeight() (int64, error)
	// Monitor 订阅最新区块，监控节点的可用性和新鲜度
	Monitor() error
	// GetBalance 获取指定地址的余额
	GetBalance(address string) (string, error)
	// GetNonce 获取指定地址的交易计数器
	GetNonce(address string) (uint64, error)
	// SubscribeToNewBlocks 订阅最新区块
	SubscribeToNewBlocks(callback func(height int64)) error
	// SubscribeToLogs 订阅新日志
	SubscribeToLogs(callback func(source, target, sender, receiver, amount string)) error
	// BuildTx 构建交易
	BuildTx(to string, value string, gasLimit uint64, gasPrice string, data []byte) (*Transaction, error)
	// EstimateGas 估算交易的Gas消耗
	EstimateGas(tx *Transaction) (uint64, error)
	// SendTx 发送交易
	SendTx(tx *Transaction) error
	// ResendTx 重新发送交易
	ResendTx(tx *Transaction) error
	// GetTransactionReceipt 获取交易收据
	GetTransactionReceipt(txHash string) (*TransactionReceipt, error)
}

type Transaction interface {
	// Hash 获取交易哈希
	Hash() string
	// RawTransaction 获取原始交易数据
	RawTransaction() ([]byte, error)
}

type TransactionReceipt interface {
	// Status 获取交易状态
	Status() bool
	// BlockNumber 获取交易所在区块号
	BlockNumber() int64
	// GasUsed 获取交易消耗的Gas
	GasUsed() uint64
}
