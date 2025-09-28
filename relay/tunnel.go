package relay

import (
	"context"
	"fmt"
	"time"

	"github.com/st-chain/me-bridge/signer"
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

	Key           *signer.Signer
	TxRecorder    *TxRecorder
	FeeCalculator *MockFeeCalculator

	msgs         chan InMsg        // 跨入消息列表（从源端订阅）
	errorHandler *BaseErrorHandler // 错误处理器

	// 控制通道
	stopCh chan struct{}
	done   chan struct{}
}

func NewInTunnel(path string, source InEndpoint, target OutEndpoint, keyManager *signer.Signer, txRecorder *TxRecorder, sequenceManager *SequenceManager, feeCalculator FeeCalculator) *InTunnel {
	return &InTunnel{
		Path:            path,
		Source:          source,
		Target:          target,
		KeyManager:      keyManager,
		TxRecorder:      txRecorder,
		SequenceManager: sequenceManager,
		FeeCalculator:   feeCalculator,
		msgs:            make(chan InMsg, 256),
		errorHandler: &BaseErrorHandler{
			Level:      LevelRelay,
			MaxRetries: 5,
			RetryDelay: time.Second * 10,
		},
		stopCh: make(chan struct{}),
		done:   make(chan struct{}),
	}
}

// Sync 从Endpoint同步Tunnel的当前跨链信息
func (t *InTunnel) Sync() {
	t.CurrSequence, t.OnChainHeight = t.Target.GetSequence()
	t.CurrNonce = t.Target.GetCurrentNonce()
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

func (t *InTunnel) Start() error {
	t.Sync()

	// 同步源端历史消息
	if err := t.GetPastMsgs(0, t.OnChainHeight); err != nil {
		return t.HandleError(context.Background(), err, map[string]interface{}{
			"operation":  "GetPastMsgs",
			"fromHeight": 0,
			"toHeight":   t.OnChainHeight,
		})
	}

	// 订阅源端入向消息
	if err := t.Source.SubscribeToInMsgs(t.msgs); err != nil {
		return t.HandleError(context.Background(), err, map[string]interface{}{
			"operation": "SubscribeToInMsgs",
		})
	}

	// 启动消息处理协程
	go t.processMessages()

	return nil
}

// processMessages 处理消息的主循环
func (t *InTunnel) processMessages() {
	defer close(t.done)

	for {
		select {
		case msg, ok := <-t.msgs:
			if !ok {
				return
			}
			t.handleMessage(msg)

		case <-t.stopCh:
			return
		}
	}
}

// handleMessage 处理单个消息
func (t *InTunnel) handleMessage(msg InMsg) {
	ctx := context.Background()
	metadata := map[string]interface{}{
		"msgNonce": msg.Nonce,
		"msgHash":  msg.TxHash,
	}

	// 处理消息
	if err := t.Target.ProcessInMsgs(make(chan InMsg, 1)); err != nil {
		// 使用分层错误处理
		if handleErr := t.HandleError(ctx, err, metadata); handleErr != nil {
			t.logError(fmt.Sprintf("Failed to handle message processing error: %v", handleErr))
		}
	}
}

// HandleError 处理 Relay 级别的错误
func (t *InTunnel) HandleError(ctx context.Context, err error, metadata map[string]interface{}) error {
	t.logError(fmt.Sprintf("Relay handling error: %v, metadata: %+v", err, metadata))

	// 使用基础错误处理器
	if handleErr := t.errorHandler.HandleError(ctx, err, metadata); handleErr != nil {
		// 错误需要升级处理，这里是最高层级，记录并可能触发告警
		t.logError(fmt.Sprintf("Fatal error at relay level: %v", handleErr))

		// 尝试恢复操作
		if err := t.recover(ctx, err, metadata); err != nil {
			return fmt.Errorf("recovery failed: %w", err)
		}
	}

	return nil
}

// recover 尝试恢复操作
func (t *InTunnel) recover(ctx context.Context, err error, metadata map[string]interface{}) error {
	t.logError("Attempting recovery at relay level")

	// 1. 重新同步状态
	t.Sync()

	// 2. 如果是消息处理错误，可能需要重新投递
	if msgNonce, ok := metadata["msgNonce"].(uint64); ok {
		// 重新构造消息并投递
		msg := InMsg{
			Nonce:  msgNonce,
			TxHash: metadata["msgHash"].(string),
			// 其他字段需要从存储中恢复
		}

		select {
		case t.msgs <- msg:
			t.logError("Message requeued for processing")
		case <-time.After(time.Second * 5):
			return fmt.Errorf("timeout requeuing message")
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Stop 停止 Tunnel
func (t *InTunnel) Stop() {
	close(t.stopCh)
	<-t.done
}

// logError 记录错误日志
func (t *InTunnel) logError(message string) {
	// TODO: 使用实际的日志系统
	// 这里简化为打印，实际应该使用结构化日志
	fmt.Printf("[ERROR] Tunnel %s: %s\n", t.Path, message)
}

// logInfo 记录信息日志
func (t *InTunnel) logInfo(message string) {
	// TODO: 使用实际的日志系统
	fmt.Printf("[INFO] Tunnel %s: %s\n", t.Path, message)
}
