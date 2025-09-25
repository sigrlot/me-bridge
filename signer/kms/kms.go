package kms

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type KMSSigner struct {
	kmsClient *kms.Client
	keyID     string
	address   common.Address
	publicKey *ecdsa.PublicKey
}

func NewKMSSigner(keyID, region string) (*KMSSigner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %v", err)
	}

	kmsClient := kms.NewFromConfig(cfg)

	// 获取公钥
	publicKey, address, err := getPublicKeyFromKMS(kmsClient, keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key from KMS: %v", err)
	}

	return &KMSSigner{
		kmsClient: kmsClient,
		keyID:     keyID,
		address:   address,
		publicKey: publicKey,
	}, nil
}

func (s *KMSSigner) GetAddress(ctx context.Context) (common.Address, error) {
	return s.address, nil
}

func (s *KMSSigner) SignData(ctx context.Context, data []byte) ([]byte, error) {
	hash := crypto.Keccak256Hash(data)
	return s.signHash(ctx, hash.Bytes())
}

func (s *KMSSigner) signHash(ctx context.Context, hash []byte) ([]byte, error) {
	input := &kms.SignInput{
		KeyId:            aws.String(s.keyID),
		Message:          hash,
		MessageType:      types.MessageTypeDigest,
		SigningAlgorithm: types.SigningAlgorithmSpecEcdsaSha256,
	}

	result, err := s.kmsClient.Sign(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("KMS signing failed: %v", err)
	}

	// 转换DER格式签名为以太坊格式
	return convertDERToEthSignature(result.Signature)
}

func (s *KMSSigner) GetPublicKey(ctx context.Context) (*ecdsa.PublicKey, error) {
	return s.publicKey, nil
}

func (s *KMSSigner) Close() error {
	return nil
}

// 辅助函数
func getPublicKeyFromKMS(client *kms.Client, keyID string) (*ecdsa.PublicKey, common.Address, error) {
	// 获取KMS密钥的公钥部分
	input := &kms.GetPublicKeyInput{
		KeyId: aws.String(keyID),
	}

	result, err := client.GetPublicKey(context.TODO(), input)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to get public key from KMS: %v", err)
	}

	// 验证密钥类型和算法
	if result.KeyUsage != types.KeyUsageTypeSignVerify {
		return nil, common.Address{}, fmt.Errorf("KMS key must have SIGN_VERIFY usage")
	}

	if result.KeySpec != types.KeySpecEccSecgP256k1 {
		return nil, common.Address{}, fmt.Errorf("KMS key must use ECC_SECG_P256K1 key spec for Ethereum")
	}

	// 解析DER格式的公钥
	publicKey, err := crypto.UnmarshalPubkey(result.PublicKey)
	if err != nil {
		return nil, common.Address{}, fmt.Errorf("failed to unmarshal public key: %v", err)
	}

	// 计算以太坊地址
	address := crypto.PubkeyToAddress(*publicKey)

	return publicKey, address, nil
}
