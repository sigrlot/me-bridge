package kms

// KMSConfig 定义 AWS KMS 签名器配置
type KMSConfig struct {
	Type   string `yaml:"type" json:"type"`     // KMS 类型，如 "aws"
	Region string `yaml:"region" json:"region"` // AWS 区域
	KeyID  string `yaml:"key_id" json:"key_id"` // KMS 密钥 ID
}
