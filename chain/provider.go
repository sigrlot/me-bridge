package chain

import "time"

// Networks stores clients for different networks
var Networks map[string]*Cluster[Client]

// ClientBuilder builds a client based on network type
func ClientBuilder(network string, config *ClientConfig) Client {
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

// ClientsBuilder creates clients for a given network configuration
func ClientsBuilder(network string, configs []*ClientConfig) []Client {
	var clients []Client
	for _, cfg := range configs {
		clients = append(clients, ClientBuilder(network, cfg))
	}

	return clients
}

func ClusterFactory(network *NetworkConfig) {
	if Networks == nil {
		Networks = make(map[string]*Cluster[Client])
	}
	clients := ClientsBuilder(network.Name, network.ClientConfigs)
	Networks[network.Name] = NewCluster[Client](clients, 30*time.Second)
}

// ClusterProvider retrieves a client by network name
func ClusterProvider(name string) *Cluster[Client] {
	if cluster, exists := Networks[name]; exists {
		return cluster
	}
	return nil
}
