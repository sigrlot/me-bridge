package relay

// Relay 代表一个跨链桥，它包含源终端、目标终端和消息队列
type Relay struct {
	Source   *InEndpoint
	Target   *OutEndpoint
	InChan   chan *InMsg  // 跨入消息列表
	BackChan chan *OutMsg // 跨出消息列表

	FeeCalculator *FeeCalculator
}

func NewRelay(source *InEndpoint, target *OutEndpoint, feeCalculator FeeCalculator) *Relay {
	return &Relay{
		Source:   source,
		Target:   target,
		InChan:   make(chan *InMsg, 1000),
		BackChan: make(chan *OutMsg, 1000),

		FeeCalculator: &feeCalculator,
	}
}

func (r *Relay) Work() error {
	// go r.Source.Work(r.InChan, r.BackChan)
	// go r.Target.Work(r.BackChan, r.InChan)

	return nil
}
