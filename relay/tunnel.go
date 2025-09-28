package relay

import "github.com/st-chain/me-bridge/signer"

type InTunnel struct {
	Path   string
	Source InEndpoint
	Target OutEndpoint

	CurrHeight   uint64 // 当前区块高度
	CurrSequence uint64 // 跨链消息序列号
	CurrNonce    uint64 // 跨链交易nonce

	KeyManager      *signer.KeyManager
	NonceManager    *NonceManager
	SequenceManager *SequenceManager
	FeeCalculator   FeeCalculator

	msgs chan InMsg // 跨入消息列表（从源端订阅）
}

func NewInTunnel(path string, source InEndpoint, target OutEndpoint, keyManager *signer.KeyManager, nonceManager *NonceManager, sequenceManager *SequenceManager, feeCalculator FeeCalculator) *InTunnel {
	return &InTunnel{
		Path:            path,
		Source:          source,
		Target:          target,
		KeyManager:      keyManager,
		NonceManager:    nonceManager,
		SequenceManager: sequenceManager,
		FeeCalculator:   feeCalculator,
		msgs:            make(chan InMsg),
	}
}

// Sync 从Endpoint同步Tunnel的当前跨链信息
func (t *InTunnel) Sync() {
	t.CurrSequence, t.CurrHeight = t.Target.GetSequence()
	t.CurrNonce = t.Target.GetNonce()
}

func (t *InTunnel) GetPastMsgs(fromHeight, toHeight uint64) error {
	msgs, err := t.Source.FilterInMsgs(fromHeight, toHeight)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		t.msgs <- msg
	}

	return nil
}

func (t *InTunnel) Start() (err error) {
	t.Sync()

	// 同步源端历史消息
	err = t.GetPastMsgs(0, t.CurrHeight)
	if err != nil {
		return err
	}

	// 订阅源端入向消息
	// TODO: 是否可以指定块高同步？
	err = t.Source.SubscribeToInMsgs(t.msgs)
	if err != nil {
		return err
	}

	// 处理入向消息
	go t.Target.ProcessInMsgs(t.msgs)
	return nil
}
