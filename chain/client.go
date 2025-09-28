package chain

import "github.com/st-chain/me-bridge/relay"

type Client interface {
	LatestHeight() (int64, error)
}

// InClient is a blockchain client that supports inbound relay operations.
type InClient interface {
	Client
	relay.InEndpoint
}

// OutClient is a blockchain client that supports outbound relay operations.
type OutClient interface {
	Client
	relay.OutEndpoint
}

// RelayLog 表示跨链日志事件
type RelayLog struct {
	TxHash   string `json:"tx_hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
	Nonce    uint64 `json:"nonce"`
}

// 实现 relay.Message 接口
func (r *RelayLog) GetTxHash() string   { return r.TxHash }
func (r *RelayLog) GetSender() string   { return r.Sender }
func (r *RelayLog) GetReceiver() string { return r.Receiver }
func (r *RelayLog) GetAmount() string   { return r.Amount }
