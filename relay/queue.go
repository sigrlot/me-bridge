package relay

// queue 实现一个简单的消息队列,对进入队列的消息进行nonce校验
type Queue[T Message] struct {
	LatestSeq uint64
	Msgs      chan T

	// stopCh chan struct{}
	// done   chan struct{}
}

func NewQueue[T Message](latestSeq uint64, bufferSize int) *Queue[T] {
	return &Queue[T]{
		LatestSeq: latestSeq,
		Msgs:      make(chan T, bufferSize),
	}
}

// Push 将消息推入队列，要求消息的nonce必须是连续的
func (q *Queue[T]) Push(msg T) {
	if msg.GetNonce() != q.LatestSeq+1 {
		// 如果nonce不连续，直接丢弃
		return
	}
	q.LatestSeq = msg.GetNonce()

	q.Msgs <- msg
}

// Pop 从队列中弹出消息，若队列为空则返回nil
func (q *Queue[T]) Pop() T {
	return <-q.Msgs
}
