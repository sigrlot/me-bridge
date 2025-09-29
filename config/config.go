package server

// APIConfig 定义 API 服务器配置
type APIConfig struct {
	IP      string `yaml:"ip" json:"ip"`           // 服务器IP地址
	Port    int32  `yaml:"port" json:"port"`       // 服务器端口
	Timeout int32  `yaml:"timeout" json:"timeout"` // 超时时间（秒）
}

// SignerConfig 定义签名器配置
type SignerConfig struct {
	Type   string `yaml:"type" json:"type"`     // 签名器类型，如 "local", "aws_kms"
	Config any    `yaml:"config" json:"config"` // 具体签名器配置
}

// KMSSignerConfig 定义 KMS 签名器配置
type KMSSignerConfig struct {
	Type   string `yaml:"type" json:"type"`     // KMS 类型，如 "aws kms"
	Region string `yaml:"region" json:"region"` // AWS 区域
	KeyID  string `yaml:"key_id" json:"key_id"` // KMS 密钥 ID
}

// NetworkConfig 定义区块链配置
type NetworkConfig struct {
	Name          string         `yaml:"network" json:"network"`               // 网络名称，如 "ethereum", "bsc", "tron"
	ChainID       string         `yaml:"chain_id" json:"chain_id"`             // 链 ID
	MaxConns      int32          `yaml:"max_conns" json:"max_conns"`           // 最大连接数
	Timeout       int64          `yaml:"timeout" json:"timeout"`               // 连接超时时间（毫秒）
	MaxRetries    int32          `yaml:"max_retries" json:"max_retries"`       // 最大重试次数
	RetryInterval int64          `yaml:"retry_interval" json:"retry_interval"` // 重试间隔（毫秒）
	ClientConfigs []ClientConfig `yaml:"target_configs" json:"target_configs"` // 目标节点配置列表
}

// ClientConfig 定义目标节点配置
type ClientConfig struct {
	Name    string `yaml:"name" json:"name"`         // 节点名称
	GRPCURL string `yaml:"grpc_url" json:"grpc_url"` // gRPC 地址
	RPCURL  string `yaml:"rpc_url" json:"rpc_url"`   // RPC 地址
	WSURL   string `yaml:"ws_url" json:"ws_url"`     // WebSocket 地址
}

// EndpointConfig 定义跨链桥端点配置
type EndpointConfig struct {
	Network         string       `yaml:"network" json:"network"`                   // 网络名称
	ConfirmBlocks   int32        `yaml:"confirm_blocks" json:"confirm_blocks"`     // 确认块数
	ContractAddress string       `yaml:"contract_address" json:"contract_address"` // 合约地址
	Signer          SignerConfig `yaml:"signer" json:"signer"`                     // 签名配置
}

// RelayConfig 定义跨链桥配置
type RelayConfig struct {
	Name   string         `yaml:"name" json:"name"`     // 桥名称
	Source EndpointConfig `yaml:"source" json:"source"` // 源端点配置
	Target EndpointConfig `yaml:"target" json:"target"` // 目标端点配置
}

// PostgresConfig 定义 PostgreSQL 数据库配置
type PostgresConfig struct {
	Host         string `yaml:"host" json:"host"`                     // 数据库主机地址
	Port         string `yaml:"port" json:"port"`                     // 数据库端口
	User         string `yaml:"user" json:"user"`                     // 用户名
	Password     string `yaml:"password" json:"password"`             // 密码
	Database     string `yaml:"database" json:"database"`             // 数据库名
	Timeout      int64  `yaml:"timeout" json:"timeout"`               // 连接超时时间（毫秒）
	MaxOpenConns int32  `yaml:"max_open_conns" json:"max_open_conns"` // 最大打开连接数
	MaxIdleConns int32  `yaml:"max_idle_conns" json:"max_idle_conns"` // 最大空闲连接数
}

// TraceConfig 定义链路追踪配置
type TraceConfig struct {
	Host       string  `yaml:"host" json:"host"`               // Jaeger 服务地址
	Enable     bool    `yaml:"enable" json:"enable"`           // 是否启用追踪
	SampleRate float32 `yaml:"sample_rate" json:"sample_rate"` // 采样率
}

type LogConfig struct {
	Level    string `yaml:"level" json:"level"`       // trace, debug, info, warn, error, fatal, panic
	Format   string `yaml:"format" json:"format"`     // json, console, ethereum
	Output   string `yaml:"output" json:"output"`     // stdout, stderr, file
	Filename string `yaml:"filename" json:"filename"` // log file name (when output is file)
}

// ServerConfig 定义服务器总配置
type ServerConfig struct {
	API      *APIConfig       `yaml:"api" json:"api"`           // API 服务器配置
	Networks []*NetworkConfig `yaml:"chains" json:"chains"`     // 区块链配置列表
	Relays   []*RelayConfig   `yaml:"bridges" json:"bridges"`   // 中继器配置列表
	Postgres *PostgresConfig  `yaml:"postgres" json:"postgres"` // PostgreSQL 配置
	Logger   *LogConfig       `yaml:"logger" json:"logger"`     // 日志配置
	Trace    *TraceConfig     `yaml:"trace" json:"trace"`       // 链路追踪配置
}
