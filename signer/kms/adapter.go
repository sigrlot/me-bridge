package kms

import (
	"encoding/asn1"
	"fmt"
	"math/big"
)

func convertDERToEthSignature(derSig []byte) ([]byte, error) {
	// DER编码的ECDSA签名结构
	type derSignature struct {
		R, S *big.Int
	}

	// 解析DER格式签名
	var sig derSignature
	if _, err := asn1.Unmarshal(derSig, &sig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DER signature: %v", err)
	}

	// 创建64字节的以太坊签名 (32字节R + 32字节S)
	// 注意：这里暂时不包含recovery ID (v)，需要在实际使用时确定
	ethSig := make([]byte, 64)

	// R值 (32字节)
	rBytes := sig.R.Bytes()
	copy(ethSig[32-len(rBytes):32], rBytes)

	// S值 (32字节)
	sBytes := sig.S.Bytes()
	copy(ethSig[64-len(sBytes):64], sBytes)

	return ethSig, nil
}
