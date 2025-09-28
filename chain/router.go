package chain

// Worker 是可排序工作单元。
type Worker interface {
	LatestHeight() int
}

// Router 包含一批Worker，并始终指向新鲜、可用的Worker。
// Router 与业务逻辑无关，仅负责在Worker之间切换。
type Router struct {
	// 当前指向的节点
	Current Worker
	// 可选的目标节点列表
	Targets []Worker `json:"targets"`
}

func NewRouter(targets []Worker) *Router {
	return &Router{
		Targets: targets,
	}
}

// Work 启动路由器，每监听所有worker的状态，并在必要时切换router执行的worker。
func (r *Router) Work() {
	// 
}

func (r *Router) Switch(newly Worker) {
	if r.Current == newly || newly == nil {
		return
	}
}
