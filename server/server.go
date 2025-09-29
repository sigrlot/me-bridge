package server

import (
	"github.com/st-chain/me-bridge/chain"
	"github.com/st-chain/me-bridge/relay"
)

type Server struct {
	
	Relays map[string]*relay.Relay
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
