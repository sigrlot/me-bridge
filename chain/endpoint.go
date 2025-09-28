package chain

import "github.com/st-chain/me-bridge/relay"

// InEndpoiont 实现 relay.InEndpoint 接口
type InEndpoiont struct {
	Config *EndpointConfig `json:"config"`

	// 终端支持的节点集合
	Clients []*relay.InEndpoint `json:"clients"`
}

func NewInEndpointWithConfig(config *EndpointConfig) *relay.InEndpoiont {
	return &relay.InEndpoiont{
		Config:  config,
		Clients: ClientsProvider(config.Network),
	}
}

// LatestHeight 返回当前终端的最新区块高度
func (e *InEndpoiont) LatestHeight() (int64, error) {
	return 0, nil
}

// Status 返回当前终端的可用性状态信息
func (e *InEndpoiont) Status() map[string]any {
	return nil
}

func (e *InEndpoiont) GetClient() *Client {
	return nil
}

// ReplaceClient 替换当前使用的节点，返回新的节点，并重新创建订阅
func (e *InEndpoiont) ReplaceClient() Client {
	return nil
}

// Monitor 监控节点的最新区块，以此判断节点的数据新鲜度和可用性
func (e *InEndpoiont) Monitor() {
}

// AutoUpdate 自动定期更新可用节点
func (e *InEndpoiont) AutoUpdate() {}

func (e *InEndpoiont) Work(local chan *RelayMsg, remote chan *RelayMsg) error {
	go e.Monitor()
	go e.AutoUpdate()

	e.GetClient().SubscribeToLogs(func(source string, target string, sender string, receiver string, amount string) {
		local <- &RelayMsg{
			Sender:   sender,
			Receiver: receiver,
			Amount:   amount,
		}
	})

	return nil
}
