package chain

import "github.com/st-chain/me-bridge/relay"

type Client interface {
	LatestHeight() (int64, error)
}

// InEndpoint 实现 relay.InEndpoint 接口
// InClient is a blockchain client that supports inbound relay operations.
// It must be able to report LatestHeight and also implement relay.InEndpoint behaviors.
type InClient interface {
	Client
	relay.InEndpoint
}

// OutEndpoint 实现 relay.OutEndpoint 接口
// OutClient is a blockchain client that supports outbound relay operations.
type OutClient interface {
	Client
	relay.OutEndpoint
}
