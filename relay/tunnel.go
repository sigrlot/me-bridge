package relay

import (
	"time"

	"github.com/st-chain/me-bridge/log"
	"github.com/st-chain/me-bridge/signer"
	"github.com/st-chain/me-bridge/types"
)

type Sequence struct {
	Height uint64
	ID     uint64
}

// InTunnel 处理跨链入向消息的通道
type InTunnel struct {
	Path     string
	Source   InEndpoint
	Target   OutEndpoint
	Sequence Sequence

	Key           signer.Signer
	Nonce         uint64 // 当前使用的nonce
	TxRecorder    *TxRecorder
	FeeCalculator *FeeCalculator

	Msgs         chan InMsg         // 跨入消息通道（从源端订阅）
	ErrorHandler types.ErrorHandler // 错误处理器
	logger       log.Logger

	// 控制通道
	// stopCh chan struct{}
	// done   chan struct{}
}

func NewInTunnel(source InEndpoint, target OutEndpoint, key signer.Signer, feeCalculator *FeeCalculator) *InTunnel {
	return &InTunnel{
		Path:          "InTunnel", // TODO: source.Name() + "->" + target.Name(),
		Source:        source,
		Target:        target,
		Key:           key,
		TxRecorder:    NewTxRecorder(0),
		FeeCalculator: feeCalculator,
		Msgs:          make(chan InMsg, 1024),
		ErrorHandler:  NewErrorHandler(3, time.Second*10),
	}
}

// Init 同步确定性跨链信息
func (t *InTunnel) Init() {
	seq, height := t.Target.GetSequence()
	t.Sequence.ID = seq
	t.Sequence.Height = height
	t.Nonce = t.Target.GetNonce(t.Key.Address())
}

// GetHistoryMsgs 获取历史跨链消息
func (t *InTunnel) GetHistoryMsgs() error {
	lastHeight := max(t.Sequence.Height, t.Source.LastHeight())
	msgs, err := t.Source.FilterInMsgs(t.Sequence.Height, lastHeight)
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		t.Msgs <- msg
	}

	return nil
}

// Start 启动 Tunnel
func (t *InTunnel) Start() error {
	t.Init()

	// 同步源端历史消息
	// TODO: 能否用订阅直接代替
	if err := t.GetHistoryMsgs(); err != nil {
		return t.HandleError(err, map[string]interface{}{
			"operation": "GetHistoryMsgs",
		})
	}

	// 订阅源端入向消息
	if err := t.Source.SubscribeToInMsgs(t.Msgs); err != nil {
		// TODO: 退出订阅
		return t.HandleError(err, map[string]interface{}{
			"operation": "SubscribeToInMsgs",
		})
	}

	// 启动消息处理协程
	if err := t.Target.ProcessInMsgs(t.Msgs); err != nil {
		// 根据错误策略处理错误
		return t.HandleError(err, map[string]interface{}{
			"operation": "ProcessInMsgs",
		})
	}

	return nil
}

// HandleError 根据错误策略处理错误
func (t *InTunnel) HandleError(err error, metadata map[string]interface{}) error {
	t.logger.Error("handling error", metadata, map[string]any{
		"error": err,
	})

	// 调用错误处理器
	if t.ErrorHandler != nil {
		t.ErrorHandler.HandleError(nil, err, metadata)
	}

	// 重启 Tunnel
	t.Stop()
	return t.Start()
}

// Stop 停止 Tunnel
func (t *InTunnel) Stop() {
	t.Source.Stop() // 停止源端订阅
	t.Target.Stop() // 停止目标端处理
}
