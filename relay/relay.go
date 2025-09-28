package relay

import "time"

// Relay 代表一个跨链桥，它包含源终端、目标终端和消息队列
type Relay struct {
	Config        *RelayConfig
	Source        InEndpoint
	Target        OutEndpoint
	NonceManager  *NonceManager
	FeeCalculator FeeCalculator

	// 消息路由的内部通道
	inChan    chan InMsg     // 跨入消息列表（从源端订阅）
	outChan   chan OutMsg    // 跨出消息列表（从目标端订阅）
	batchChan chan *BatchMsg // 批量消息列表（从目标端订阅）

	// 控制通道
	stopCh chan struct{}
	done   chan struct{}
}

func NewRelay(config *RelayConfig, source InEndpoint, target OutEndpoint, feeCalculator FeeCalculator, startNonce uint64) *Relay {
	return &Relay{
		Config:        config,
		Source:        source,
		Target:        target,
		NonceManager:  NewNonceManager(startNonce),
		FeeCalculator: feeCalculator,

		inChan:    make(chan InMsg, 1000),
		outChan:   make(chan OutMsg, 1000),
		batchChan: make(chan *BatchMsg, 1000),

		stopCh: make(chan struct{}),
		done:   make(chan struct{}),
	}
}

func (r *Relay) Work() error {
	// 1. 订阅源端入向消息
	sourceInChan, err := r.Source.SubscribeToInMsgs()
	if err != nil {
		return err
	}

	// 2. 订阅目标端出向消息
	targetOutChan, err := r.Target.SubscribeToOutMsgs()
	if err != nil {
		return err
	}

	// 3. 订阅目标端批量消息
	targetBatchChan, err := r.Target.SubscribeToBatchMsgs()
	if err != nil {
		return err
	}

	// 启动消息处理协程
	go r.processInboundMessages(sourceInChan)
	go r.processOutboundMessages(targetOutChan)
	go r.processBatchMessages(targetBatchChan)
	go r.processTargetMessages()
	go r.processSourceMessages()
	go r.cleanupRoutine()

	return nil
}

// processInboundMessages 处理来自源链的消息（桥接入向消息）
func (r *Relay) processInboundMessages(inChan <-chan InMsg) {
	for {
		select {
		case msg := <-inChan:
			r.inChan <- msg
		case <-r.stopCh:
			return
		}
	}
}

// processOutboundMessages 处理来自目标链的确认消息
func (r *Relay) processOutboundMessages(outChan <-chan OutMsg) {
	for {
		select {
		case msg := <-outChan:
			r.outChan <- msg
		case <-r.stopCh:
			return
		}
	}
}

// processBatchMessages 处理来自目标链的批量消息
func (r *Relay) processBatchMessages(batchChan <-chan *BatchMsg) {
	for {
		select {
		case msg := <-batchChan:
			r.batchChan <- msg
		case <-r.stopCh:
			return
		}
	}
}

// processTargetMessages 处理入向消息并发送到目标链
func (r *Relay) processTargetMessages() {
	for {
		select {
		case inMsg := <-r.inChan:
			// 将InMsg转换为OutMsg并进行nonce管理
			outMsg := &OutMsg{
				Sender:   inMsg.Sender,
				Receiver: inMsg.Receiver,
				Amount:   inMsg.Amount,
			}

			// 计算费用
			if r.FeeCalculator != nil {
				// 费用计算逻辑在此处
				// 目前只是直接通过
			}

			// 为异步交易分配nonce
			nonce := r.NonceManager.AllocateNonce(outMsg)

			// 发送到目标链进行处理
			go func(msg *OutMsg, allocatedNonce uint64) {
				// 这将由目标端点处理
				// 目标端将通过MarkSubmitted回调交易哈希
				if err := r.Target.ProcessInMsgs(make(chan InMsg)); err != nil {
					r.NonceManager.MarkFailed(allocatedNonce)
				}
			}(outMsg, nonce)

		case <-r.stopCh:
			return
		}
	}
}

// processSourceMessages 处理出向确认消息并发送到源链
func (r *Relay) processSourceMessages() {
	for {
		select {
		case outMsg := <-r.outChan:
			// 处理来自目标链的确认消息
			go func(msg OutMsg) {
				if err := r.Source.ProcessOutMsgs(make(chan OutMsg)); err != nil {
					// 处理错误
				}
			}(outMsg)

		case <-r.stopCh:
			return
		}
	}
}

// cleanupRoutine 定期清理过期交易
func (r *Relay) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 清理超过1小时的交易
			r.NonceManager.CleanupStale(time.Hour)

			// 重试失败的交易
			retryable := r.NonceManager.GetRetryableTxs(3)
			for _, tx := range retryable {
				// 重试逻辑在此处
				_ = tx // 占位符
			}

		case <-r.stopCh:
			return
		}
	}
}

// Stop 优雅地停止中继器
func (r *Relay) Stop() {
	close(r.stopCh)
	<-r.done
}

// GetStatus 返回中继器的当前状态
func (r *Relay) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"current_nonce":  r.NonceManager.GetCurrentNonce(),
		"pending_count":  r.NonceManager.GetPendingCount(),
		"in_chan_len":    len(r.inChan),
		"out_chan_len":   len(r.outChan),
		"batch_chan_len": len(r.batchChan),
	}
}
