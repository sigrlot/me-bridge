package relay

// queue 实现一个简单的消息队列,对进入队列的消息进行nonce校验
type Queue struct {
	LatestSeq uint64
	msgs      chan Message
}

func NewQueue(latestSeq uint64, bufferSize int) *Queue {
	return &Queue{
		LatestSeq: latestSeq,
		msgs:      make(chan Message, bufferSize),
	}
}

// Push 将消息推入队列，要求消息的nonce必须是连续的
func (q *Queue) Push(msg Message) {
	if msg.GetNonce() != q.LatestSeq+1 {
		// 如果nonce不连续，直接丢弃
		return
	}
	q.LatestSeq = msg.GetNonce()

	q.msgs <- msg
}

// Pop 从队列中弹出消息，若队列为空则返回nil
func (q *Queue) Pop() Message {
	return <-q.msgs
}
