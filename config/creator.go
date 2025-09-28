package server

import (
	"github.com/st-chain/me-bridge/chain"
	"github.com/st-chain/me-bridge/relay"
)

func NewRelayWithConfig(config *relay.RelayConfig) *relay.Relay {
	source := NewInEndpointWithConfig(config.Source)
	target := NewOutEndpointWithConfig(config.Target)

	return &Relay{
		Source:   source,
		Target:   target,
		InChan:   make(chan *Message, 1000),
		BackChan: make(chan *Message, 1000),

		FeeCalculator: &FeeCalculator{},
	}
}

func NewInEndpointWithConfig(config *EndpointConfig) *relay.InEndpoiont {
	return &relay.InEndpoiont{
		Config:  config,
		Clients: ClientsProvider(config.Network),
	}
}

