package relay

// Relay 代表一个跨链桥，它包含源终端、目标终端和消息队列
type Relay struct {
	Source   *Endpoint
	Target   *Endpoint
	InChan   chan *Message // 跨入消息列表
	BackChan chan *Message // 跨出消息列表

	FeeCalculator *FeeCalculator
}

func NewRelay(source *Endpoint, target *Endpoint, feeCalculator FeeCalculator) *Relay {
	return &Relay{
		Source:   source,
		Target:   target,
		InChan:   make(chan *Message, 1000),
		BackChan: make(chan *Message, 1000),

		FeeCalculator: &feeCalculator,
	}
}

func NewRelayWithConfig(config *RelayConfig) *Relay {
	source := NewEndpointWithConfig(config.Source)
	target := NewEndpointWithConfig(config.Target)

	return &Relay{
		Source:   source,
		Target:   target,
		InChan:   make(chan *Message, 1000),
		BackChan: make(chan *Message, 1000),

		FeeCalculator: &FeeCalculator{},
	}
}

func (r *Relay) Work() error {
	go r.Source.Work(r.InChan, r.BackChan)
	go r.Target.Work(r.BackChan, r.InChan)

	return nil
}
