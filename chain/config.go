package chain

// NetworkConfig 定义区块链配置
type NetworkConfig struct {
	Name          string          `yaml:"network" json:"network"`               // 网络名称，如 "ethereum", "bsc", "tron"
	ChainID       string          `yaml:"chain_id" json:"chain_id"`             // 链 ID
	MaxConns      int32           `yaml:"max_conns" json:"max_conns"`           // 最大连接数
	ClientConfigs []*ClientConfig `yaml:"target_configs" json:"target_configs"` // 目标节点配置列表
}

// ClientConfig 定义目标节点配置
type ClientConfig struct {
	Name     string `yaml:"name" json:"name"`           // 节点名称
	URI      string `yaml:"uri" json:"uri"`             // 节点 URI
	GRPCPort string `yaml:"grpc_port" json:"grpc_port"` // gRPC 端口
	RPCPort  string `yaml:"rpc_port" json:"rpc_port"`   // RPC 端口
	WSPort   string `yaml:"ws_port" json:"ws_port"`     // WebSocket 端口
}
