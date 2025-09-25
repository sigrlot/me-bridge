package relay

// Relay 代表一个跨链桥，它包含源终端、目标终端和消息队列
type Relay struct {
	Source   *Endpoint
	Target   *Endpoint
	InQueue  *Queue // 跨入消息列表
	OutQueue *Queue // 跨出消息列表

	FeeCalculator *FeeCalculator
}

func NewRelay(source *Endpoint, target *Endpoint, feeCalculator FeeCalculator) *Relay {
	return &Relay{
		Source:   source,
		Target:   target,
		InQueue:  NewQueue(source, target),
		OutQueue: NewQueue(source, target),
	}
}

func NewRelayWithConfig(config *RelayConfig) *Relay {
	source := NewEndpointWithConfig(config.Source)
	target := NewEndpointWithConfig(config.Target)

	return &Relay{
		Source:   source,
		Target:   target,
		InQueue:  NewQueue(source, target),
		OutQueue: NewQueue(target, source),
	}
}

func (b *Relay) Start() error {
	log.Info("start relay")
	// 启动监控，定期更新节点列表
	go b.Source.Monitor()
	go b.Target.Monitor()

	go b.Source.Subscribe(b.InQueue)
	go b.Target.Subscribe(b.OutQueue)

	// 启动工作协程，处理消息队列
	go b.InQueue.Worker(b.Target)
	go b.OutQueue.Worker(b.Source)
}
