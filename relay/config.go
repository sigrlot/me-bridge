package relay

import logger "github.com/st-chain/me-bridge/log"

var log *logger.Logger

func InitLogger() {
	log = logger.GetLogger().WithComponent("relay")
}

// RelayConfig 定义跨链桥配置
type RelayConfig struct {
	Name          string         `yaml:"name" json:"name"`                     // 桥名称
	Source        *EndpointConfig `yaml:"source" json:"source"`                 // 源端点配置
	Target        *EndpointConfig `yaml:"target" json:"target"`                 // 目标端点配置
	MaxRetries    int32          `yaml:"max_retries" json:"max_retries"`       // 最大重试次数
	RetryInterval int64          `yaml:"retry_interval" json:"retry_interval"` // 重试间隔（毫秒）
}

// EndpointConfig 定义端点配置
type EndpointConfig struct {
	Network         string `yaml:"network" json:"network"`                   // 网络名称
	ContractAddress string `yaml:"contract_address" json:"contract_address"` // 合约地址
	SourceKeyID     string `yaml:"source_key_id" json:"source_key_id"`       // 签名密钥 ID
	ConfirmBlocks   int32  `yaml:"confirm_blocks" json:"confirm_blocks"`     // 确认块数
}
