package bsc

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// BSCClient is a client for interacting with the Binance Smart Chain (BSC) network.
type BSCClient struct {
	client     *ethclient.Client
	rpcClient  *rpc.Client
	chainID    *big.Int
	timeout    time.Duration
	privateKey *ecdsa.PrivateKey
}

// Config holds BSC client configuration
type Config struct {
	RPCURL     string
	WSUrl      string
	ChainID    int64
	Timeout    time.Duration
	PrivateKey string
}

// NewBSCClient creates a new BSC client instance
func NewBSCClient(config *Config) (*BSCClient, error) {
	// Connect to BSC node
	client, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, err
	}

	// Connect RPC client for advanced operations
	rpcClient, err := rpc.Dial(config.RPCURL)
	if err != nil {
		return nil, err
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &BSCClient{
		client:     client,
		rpcClient:  rpcClient,
		chainID:    big.NewInt(config.ChainID),
		timeout:    timeout,
		privateKey: privateKey,
	}, nil
}

// GetBalance returns the balance of an address
func (c *BSCClient) GetBalance(address common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.BalanceAt(ctx, address, nil)
}

// GetNonce returns the nonce for an address
func (c *BSCClient) GetNonce(address common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.PendingNonceAt(ctx, address)
}

// GetGasPrice returns current gas price
func (c *BSCClient) GetGasPrice() (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.SuggestGasPrice(ctx)
}

// GetLatestBlock returns the latest block number
func (c *BSCClient) GetLatestBlock() (*types.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.HeaderByNumber(ctx, nil)
}

// GetTransactionReceipt returns transaction receipt
func (c *BSCClient) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.TransactionReceipt(ctx, txHash)
}

// Close closes the client connections
func (c *BSCClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}
