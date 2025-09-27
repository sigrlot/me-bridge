package chain

// Networks stores clients for different networks
var Networks map[string][]*Monitor

// ClientBuilder builds a client based on network type
func ClientBuilder(network string, config *ClientConfig) *Monitor {
	switch network {
	case "ethereum":
		return NewEthereumClient(config)
	case "binance":
		return NewBSCClient(config)
	case "tron":
		return NewTronClient(config)
	}
	panic("unsupported network")
}

// ClientsFactory creates clients for a given network configuration
func ClientsFactory(network string, configs []*ClientConfig) []*Monitor {
	var clients []*Monitor
	for _, cfg := range configs {
		clients = append(clients, ClientBuilder(network, cfg))
	}
	// add clients to client pool
	Networks[network] = clients

	return clients
}

// ClientsProvider retrieves a client by network name
func ClientsProvider(name string) []*Monitor {
	if clients, exists := Networks[name]; exists && len(clients) > 0 {
		return clients
	}
	return nil
}
