package signer

import (
	"context"
)

// Signer 定义了签名器的通用接口
type Signer interface {
	// Address 获取签名器的以太坊地址
	Address() string

	// PublicKey 获取公钥
	PublicKey() string

	// SignData 签名任意数据
	SignData(ctx context.Context, data []byte) ([]byte, error)

	// Close 关闭签名器并清理资源
	Close() error
}
