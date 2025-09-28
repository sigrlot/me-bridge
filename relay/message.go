package relay

// Message 通用消息接口
type Message interface {
	GetNonce() uint64
}

type InMsg struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type OutMsg struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type BatchMsg struct {
	Nonces uint64 `json:"nonces"`
	Data   []byte `json:"data"`
}
