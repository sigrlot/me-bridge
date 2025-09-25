package server

import (
	"github.com/st-chain/me-bridge/api"
	"github.com/st-chain/me-bridge/chain"
	"github.com/st-chain/me-bridge/db"
	"github.com/st-chain/me-bridge/log"
	"github.com/st-chain/me-bridge/relay"
	"github.com/st-chain/me-bridge/signer/kms"
	"github.com/st-chain/me-bridge/trace"
)

// ServerConfig 定义服务器总配置
type ServerConfig struct {
	API      *api.APIConfig         `yaml:"api" json:"api"`           // API 服务器配置
	KMS      *kms.KMSConfig         `yaml:"kms" json:"kms"`           // KMS 配置
	Networks []*chain.NetworkConfig `yaml:"chains" json:"chains"`     // 区块链配置列表
	Relays    []*relay.RelayConfig   `yaml:"bridges" json:"bridges"`   // 中继器配置列表
	Postgres *db.PostgresConfig     `yaml:"postgres" json:"postgres"` // PostgreSQL 配置
	Logger   *log.LoggerConfig      `yaml:"logger" json:"logger"`     // 日志配置
	Trace    *trace.TraceConfig     `yaml:"trace" json:"trace"`       // 链路追踪配置
}
