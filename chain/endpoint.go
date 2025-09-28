package chain

import (
	"time"

	"github.com/st-chain/me-bridge/relay"
)


// InEndpoint 实现 relay.InEndpoint 接口，通过 Cluster 统一管理多个节点。
type InEndpoint struct {
	Config  *relay.EndpointConfig `json:"config"`
	cluster *Cluster[InClient]
}

// LatestHeight 返回当前终端的最新区块高度
func (e *InEndpoint) LatestHeight() (int64, error) { return e.cluster.Current().LatestHeight() }

// Status 返回当前终端的可用性状态信息
func (e *InEndpoint) Status() map[string]any { return nil }

func (e *InEndpoint) GetClient() InClient { return e.cluster.Current() }

// ReplaceClient 替换当前使用的节点，返回新的节点，并重新创建订阅
func (e *InEndpoint) ReplaceClient() InClient { return e.cluster.ReplaceClient() }

// Monitor 监控节点的最新区块，以此判断节点的数据新鲜度和可用性
func (e *InEndpoint) Monitor() { e.cluster.Start() }

// AutoUpdate 自动定期更新可用节点
func (e *InEndpoint) AutoUpdate() {}

// Implement relay.InEndpoint by delegating to the current client in cluster.
func (e *InEndpoint) SubscribeToInMsgs(address string) (<-chan relay.InMsg, error) {
	return e.GetClient().SubscribeToInMsgs(address)
}

func (e *InEndpoint) ProcessOutMsgs(msgs <-chan relay.OutMsg) error {
	return e.GetClient().ProcessOutMsgs(msgs)
}

// NewInEndpoint constructs an InEndpoint with clients and monitoring interval.
func NewInEndpoint(config *relay.EndpointConfig, clients []InClient, monitorInterval time.Duration) *InEndpoint {
	ep := &InEndpoint{
		Config:  config,
		cluster: NewCluster[InClient](clients, monitorInterval),
	}
	ep.cluster.Start()
	return ep
}

// OutEndpoint 实现 relay.OutEndpoint 接口，通过 Cluster 统一管理多个节点。
type OutEndpoint struct {
	Config  *relay.EndpointConfig `json:"config"`
	cluster *Cluster[OutClient]
}

func (e *OutEndpoint) LatestHeight() (int64, error) { return e.cluster.Current().LatestHeight() }

func (e *OutEndpoint) GetClient() OutClient { return e.cluster.Current() }

func (e *OutEndpoint) ReplaceClient() OutClient { return e.cluster.ReplaceClient() }

// Implement relay.OutEndpoint by delegating to current client
func (e *OutEndpoint) ProcessInMsgs(msgs <-chan relay.InMsg) error {
	return e.GetClient().ProcessInMsgs(msgs)
}

func (e *OutEndpoint) SubscribeToOutMsgs() (<-chan relay.OutMsg, error) {
	return e.GetClient().SubscribeToOutMsgs()
}

func (e *OutEndpoint) ConfirmOutMsgs(msgs []relay.OutMsg) error {
	return e.GetClient().ConfirmOutMsgs(msgs)
}

func (e *OutEndpoint) SubscribeToBatchMsgs() (<-chan relay.BatchMsg, error) {
	return e.GetClient().SubscribeToBatchMsgs()
}

// NewOutEndpoint constructs an OutEndpoint with clients and monitoring interval.
func NewOutEndpoint(config *relay.EndpointConfig, clients []OutClient, monitorInterval time.Duration) *OutEndpoint {
	ep := &OutEndpoint{
		Config:  config,
		cluster: NewCluster[OutClient](clients, monitorInterval),
	}
	ep.cluster.Start()
	return ep
}
