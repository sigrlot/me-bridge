package bsc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	bscclient "github.com/ethereum/go-ethereum/ethclient"

	"github.com/st-chain/me-bridge/chain"
	"github.com/st-chain/me-bridge/log"
	"github.com/st-chain/me-bridge/relay"
)

var _ relay.Client = (*Client)(nil)

// Client is a client for interacting with the Binance Smart Chain (BSC) network.
type Client struct {
	Network *chain.NetworkConfig
	Config  *chain.ClientConfig

	latestHeight uint64

	Client   *bscclient.Client
	WsClient *bscclient.Client
	// key    signer.Signer
	logger *log.Logger
}

// NewClient creates a new BSC client instance
func NewClient(network *chain.NetworkConfig, config *chain.ClientConfig) (*Client, error) {
	// Connect to BSC node
	client, err := bscclient.Dial(config.RPCURL)
	if err != nil {
		return nil, err
	}

	wsClient, err := ethclient.Dial(config.WSURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		Network: network,
		Config:  config,

		Client:   client,
		WsClient: wsClient,
		logger:   log.WithComponent("bsc-client"),
	}, nil
}

func (c *Client) LatestHeight() uint64 {
	return c.latestHeight
}

// GetBalance returns the balance of an address
func (c *Client) GetBalance(address common.Address, blockNumber *big.Int) (*big.Int, error) {
	return c.Client.BalanceAt(context.Background(), address, blockNumber)
}

// GetNonce returns the nonce for an address
func (c *Client) GetNonce(address common.Address) (uint64, error) {
	return c.Client.PendingNonceAt(context.Background(), address)
}

// GetGasPrice returns current gas price
func (c *Client) GetGasPrice() (*big.Int, error) {
	return c.Client.SuggestGasPrice(context.Background())
}

// GetLatestBlock returns the latest block number
func (c *Client) GetLatestHeight() (uint64, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

// GetTransactionReceipt returns transaction receipt
func (c *Client) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return c.Client.TransactionReceipt(context.Background(), txHash)
}

// Close closes the client connections
func (c *Client) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
	if c.WsClient != nil {
		c.WsClient.Close()
	}
}
