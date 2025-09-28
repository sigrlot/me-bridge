package relay

type Queue struct {
	Height uint64
	Msgs   []byte
}


func (q *Queue) Push(msg []byte) {
	
}