package relay

type InMsg struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

type OutMsg struct {
	Nonce uint64 `json:"nonce"`
	Data  []byte `json:"data"`
}
