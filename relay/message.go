package relay

type InMsg struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type OutMsg struct {
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type BatchMsg struct {
	Nonces uint64 `json:"nonces"`
	Data   []byte `json:"data"`
}
