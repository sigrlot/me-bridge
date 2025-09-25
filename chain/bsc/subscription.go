package bsc

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	EthTopic    = "eth_topic"    // 以太坊订阅内容
	TronTopic   = "tron_topic"   // 波场订阅内容
	CosmosTopic = "cosmos_topic" // 宇宙订阅内容
)


// SubscriptionManager handles event subscriptions
type SubscriptionManager struct {
	client    *ethclient.Client
	wsClient  *ethclient.Client
	callbacks map[string]func(inter	TronTopic   = "tron_topic"   // 波场订阅内容
face{})
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(rpcURL, wsURL string) (*SubscriptionManager, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	wsClient, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, err
	}

	return &SubscriptionManager{
		client:    client,
		wsClient:  wsClient,
		callbacks: make(map[string]func(any)),
	}, nil
}

// SubscribeToNewBlocks subscribes to new block headers
func (sm *SubscriptionManager) SubscribeToNewBlocks(callback func(*types.Header)) error {
	headers := make(chan *types.Header)

	sub, err := sm.wsClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Subscription error: %v", err)
				return
			case header := <-headers:
				callback(header)
			}
		}
	}()

	return nil
}

// SubscribeToLogs subscribes to contract logs
func (sm *SubscriptionManager) SubscribeToLogs(addresses []common.Address, topics [][]common.Hash, callback func(types.Log)) error {
	query := ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
	}

	logs := make(chan types.Log)
	sub, err := sm.wsClient.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				log.Printf("Log subscription error: %v", err)
				return
			case vLog := <-logs:
				callback(vLog)
			}
		}
	}()

	return nil
}

// GetLogs retrieves logs for a specific filter
func (sm *SubscriptionManager) GetLogs(fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
		Topics:    topics,
	}

	return sm.client.FilterLogs(context.Background(), query)
}

// Close closes all connections
func (sm *SubscriptionManager) Close() {
	if sm.client != nil {
		sm.client.Close()
	}
	if sm.wsClient != nil {
		sm.wsClient.Close()
	}
}
