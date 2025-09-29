package relay

// Message 通用消息接口
type Message interface {
	GetNonce() uint64
}

// InMsg 代表从源端接收到的跨入消息
type InMsg struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (m InMsg) GetNonce() uint64 {
	return m.Nonce
}

// OutMsg 代表从目标端发送的跨出消息
type OutMsg struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (m OutMsg) GetNonce() uint64 {
	return m.Nonce
}

// BatchMsg 代表从目标端发送的预签名批量消息
type BatchMsg struct {
	Nonces uint64 `json:"nonces"`
	Data   []byte `json:"data"`
}

func (m BatchMsg) GetNonce() uint64 {
	return m.Nonces
}
