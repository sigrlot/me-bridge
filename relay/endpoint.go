package relay

import (
	"github.com/st-chain/me-bridge/chain"
)

// Endpoint 代表跨链桥的一端，并提供跨链桥相关的封装方法。
// Endpoint 包含一批节点，以避免节点的单点问题，始终保证跨链桥始终使用新鲜、可用的节点。
type Endpoint struct {
	Config *EndpointConfig `json:"config"`

	// 终端支持的节点集合
	Clients []*chain.Client `json:"clients"`
}

func NewEndpointWithConfig(config *EndpointConfig) *Endpoint {
	return &Endpoint{
		Config:  config,
		Clients: chain.ClientsProvider(config.Network),
	}
}

// LatestHeight 返回当前终端的最新区块高度
func (e *Endpoint) LatestHeight() (int64, error) {
	return 0, nil
}

// Status 返回当前终端的可用性状态信息
func (e *Endpoint) Status() map[string]any {
	return nil
}

func (e *Endpoint) GetClient() chain.Client {
	return nil
}

// ReplaceClient 替换当前使用的节点，返回新的节点，并重新创建订阅
func (e *Endpoint) ReplaceClient() chain.Client {
	return nil
}

// Monitor 监控节点的最新区块，以此判断节点的数据新鲜度和可用性
func (e *Endpoint) Monitor() {
}

// AutoUpdate 自动定期更新可用节点
func (e *Endpoint) AutoUpdate() {}

func (e *Endpoint) Start() error {
	go e.Monitor()
	go e.AutoUpdate()

	return nil
}
