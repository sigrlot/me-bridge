package signer

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Signer 定义了签名器的通用接口
type Signer interface {
	// GetAddress 获取签名器的以太坊地址
	GetAddress(ctx context.Context) (common.Address, error)

	// SignTransaction 签名以太坊交易
	SignTransaction(ctx context.Context, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	// SignData 签名任意数据
	SignData(ctx context.Context, data []byte) ([]byte, error)

	// GetPublicKey 获取公钥
	GetPublicKey(ctx context.Context) (*ecdsa.PublicKey, error)

	// Close 关闭签名器并清理资源
	Close() error

	// HandleError 处理签名器级别的错误
	HandleError(ctx context.Context, err error, metadata map[string]interface{}) error
}
