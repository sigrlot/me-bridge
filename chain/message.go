package chain

import "fmt"

type RelayLog struct {
	Nonce    uint64 `json:"nonce"`
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (m *RelayLog) GetTxHash() string {
	return m.TxHash
}

func (m *RelayLog) GetSender() string {
	return m.Sender
}

func (m *RelayLog) GetReceiver() string {
	return m.Receiver
}

func (m *RelayLog) GetAmount() string {
	return m.Amount
}

func (m *RelayLog) String() string {
	return fmt.Sprintf("RelayMsg{TxHash: %s, Sender: %s, Receiver: %s, Amount: %s}", m.TxHash, m.Sender, m.Receiver, m.Amount)
}
