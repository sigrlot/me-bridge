package chain

// EndpointConfig 定义端点配置
type EndpointConfig struct {
	Network         string `yaml:"network" json:"network"`                   // 网络名称
	ContractAddress string `yaml:"contract_address" json:"contract_address"` // 合约地址
	SourceKeyID     string `yaml:"source_key_id" json:"source_key_id"`       // 签名密钥 ID
	ConfirmBlocks   int32  `yaml:"confirm_blocks" json:"confirm_blocks"`     // 确认块数
}

// NetworkConfig 定义区块链配置
type NetworkConfig struct {
	Name          string          `yaml:"network" json:"network"`               // 网络名称，如 "ethereum", "bsc", "tron"
	ChainID       string          `yaml:"chain_id" json:"chain_id"`             // 链 ID
	ClientConfigs []*ClientConfig `yaml:"target_configs" json:"target_configs"` // 目标节点配置列表
	MaxConns      int32           `yaml:"max_conns" json:"max_conns"`           // 最大连接数
	Timeout       int64           `yaml:"timeout" json:"timeout"`               // 超时时间，单位秒
	// RetryCount    int             `yaml:"retry_count" json:"retry_count"`       // 重试次数
	// RetryDelay    int64           `yaml:"retry_delay" json:"retry_delay"`       // 重试延迟，单位秒
}

// ClientConfig 定义目标节点配置
type ClientConfig struct {
	Name    string `yaml:"name" json:"name"`         // 节点名称
	GRPCURL string `yaml:"grpc_url" json:"grpc_url"` // gRPC 地址
	RPCURL  string `yaml:"rpc_url" json:"rpc_url"`   // RPC 地址
	WSURL   string `yaml:"ws_url" json:"ws_url"`     // WebSocket 地址
}
