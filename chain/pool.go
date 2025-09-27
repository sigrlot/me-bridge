package chain

// ClientPool 观察同网络下的多个节点的最新状态，提供可用的节点，避免节点单点问题。
type ClientPool struct {
	Clients []*Client `json:"clients"`
}

func NewClientPool(clients []*Client) *ClientPool {
	return &ClientPool{
		Clients: clients,
	}
}

// Status 返回当前终端的可用性状态信息
func (p *ClientPool) Status() map[string]any {
	return nil
}

// LatestHeight 返回当前终端的最新区块高度
func (p *ClientPool) LatestHeight() (int64, error) {
	return 0, nil
}

func (p *ClientPool) GetClient() *Client {
	return nil
}

// ReplaceClient 替换当前使用的节点，返回新的节点，并重新创建订阅
func (p *ClientPool) ReplaceClient() *Client {
	return nil
}

// Monitor 监控节点的最新区块，以此判断节点的数据新鲜度和可用性
func (p *ClientPool) Monitor() {
}

// AutoUpdate 自动定期更新可用节点
func (p *ClientPool) AutoUpdate() {}

func (p *ClientPool) Work() error {
	go p.Monitor()
	go p.AutoUpdate()

	return nil
}
