package server

import (
	"github.com/st-chain/me-bridge/chain"
	"github.com/st-chain/me-bridge/relay"
)

type Server struct {
	config *ServerConfig
	Relays map[string]*relay.Relay
}

func NewServerWithConfig(config *ServerConfig) *Server {
	for _, netConfig := range config.Networks {
		chain.ClientsFactory(netConfig.Name, netConfig.ClientConfigs)
	}

	relays := make(map[string]*relay.Relay)
	for _, relayConfig := range config.Relays {
		relay := relay.NewRelayWithConfig(relayConfig)
		relays[relayConfig.Name] = relay
	}

	return &Server{
		config: config,
		Relays: relays,
	}
}

func (s *Server) Start() error {
	for _, relay := range s.Relays {
		if err := relay.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Stop() error {
	for _, relay := range s.Relays {
		if err := relay.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Status() string {
	return "running"
}

func (s *Server) Restart() error {
	return nil
}
