package chain

import "github.com/st-chain/me-bridge/relay"

type Client interface {
	LatestHeight() (int64, error)
}

// InClient is a blockchain client that supports inbound relay operations.
type InClient interface {
	Client
	relay.InEndpoint
}

// OutClient is a blockchain client that supports outbound relay operations.
type OutClient interface {
	Client
	relay.OutEndpoint
}
