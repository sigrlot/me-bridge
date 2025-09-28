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

// 通过委托给集群中的当前客户端来实现 relay.InEndpoint
func (e *InEndpoint) SubscribeToInMsgs(contract string, c <-chan *relay.InMsg) error {
	return e.GetClient().SubscribeToInMsgs(contract, c)
}

func (e *InEndpoint) ProcessOutMsgs(msgs <-chan *relay.OutMsg) error {
	return e.GetClient().ProcessOutMsgs(msgs)
}

// NewInEndpoint 使用客户端和监控间隔构建 InEndpoint
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

// 通过委托给当前客户端来实现 relay.OutEndpoint
func (e *OutEndpoint) ProcessInMsgs(msgs <-chan *relay.InMsg) error {
	return e.GetClient().ProcessInMsgs(msgs)
}

func (e *OutEndpoint) SubscribeToOutMsgs(msgs <-chan *relay.OutMsg) error {
	return e.GetClient().SubscribeToOutMsgs(msgs)
}

func (e *OutEndpoint) ConfirmOutMsgs(msgs []relay.OutMsg) error {
	return e.GetClient().ConfirmOutMsgs(msgs)
}

func (e *OutEndpoint) SubscribeToBatchMsgs() (<-chan relay.BatchMsg, error) {
	return e.GetClient().SubscribeToBatchMsgs()
}

// NewOutEndpoint 使用客户端和监控间隔构建 OutEndpoint
func NewOutEndpoint(config *relay.EndpointConfig, clients []OutClient, monitorInterval time.Duration) *OutEndpoint {
	ep := &OutEndpoint{
		Config:  config,
		cluster: NewCluster[OutClient](clients, monitorInterval),
	}
	ep.cluster.Start()
	return ep
}
