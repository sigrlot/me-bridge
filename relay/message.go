package relay

import "fmt"

type RelayMsg struct {
	TxHash   string `json:"hash"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   string `json:"amount"`
}

func (m *RelayMsg) String() string {
	return fmt.Sprintf("RelayMsg{TxHash: %s, Sender: %s, Receiver: %s, Amount: %s}", m.TxHash, m.Sender, m.Receiver, m.Amount)
}

func (m *RelayMsg) Map() map[string]any {
	return map[string]any{
		"hash":     m.TxHash,
		"sender":   m.Sender,
		"receiver": m.Receiver,
		"amount":   m.Amount,
	}
}
