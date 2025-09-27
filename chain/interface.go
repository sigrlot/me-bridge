package chain

// Monitor 封装了区块链client的差异性，提供了统一的跨链桥相关方法
type Monitor interface {
	// LatestHeight 获取最新区块高度
	LatestHeight() uint64RelayLog
	// Monitor 订阅最新区块，监控节点的可用性和新鲜度
	TrackHeight() error
}
