package relay

type Queue struct {
	Source *Endpoint `json:"source"`
	Target *Endpoint `json:"target"`

	Msgs     chan *RelayMsg `json:"-"`
	WaitMsgs chan *RelayMsg `json:"-"`
}

func NewQueue(source *Endpoint, target *Endpoint) *Queue {
	return &Queue{
		Source:   source,
		Target:   target,
		Msgs:     make(chan *RelayMsg, 100),
		WaitMsgs: make(chan *RelayMsg, 100),
	}
}

func (pc *Queue) Close() {
	close(pc.Msgs)
}

func (pc *Queue) InQueue(msg *RelayMsg) {
	pc.Msgs <- msg
}

func (pc *Queue) OutQueue() *RelayMsg {
	return <-pc.Msgs
}

func (pc *Queue) Worker(target *Endpoint) {
	for msg := range pc.Msgs {
		// 处理消息
		_ = msg
	}
}
