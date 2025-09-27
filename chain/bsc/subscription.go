package bsc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/st-chain/me-bridge/chain"
	"github.com/st-chain/me-bridge/relay"
)

// BSC 跨链事件订阅主题
var RelayTopic = [][]common.Hash{
	{common.HexToHash("")},
}

// TrackHeight tracks the latest block height
func (c *Client) TrackHeight() error {
	c.logger.Debug("Starting to track latest block height")

	headers := make(chan *types.Header)
	sub, err := c.WsClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		c.logger.Error("Failed to subscribe to new block", map[string]any{
			"wsurl": c.Config.WSURL,
			"error": err,
		})
		return err
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				c.logger.Error("Subscription error", map[string]any{
					"wsurl": c.Config.WSURL,
					"error": err,
				})
				return
			case header := <-headers:
				c.latestHeight = header.Number.Uint64()
			}
		}
	}()

	return nil
}

func (c *Client) ToRelayLog(vLog types.Log) (*chain.RelayLog, error) {
	relayLog := &chain.RelayLog{
		TxHash:   vLog.TxHash.Hex(),
		Sender:   "", // 需要解析Data字段获取
		Receiver: "", // 需要解析Data字段获取
		Amount:   "", // 需要解析Data字段获取
	}
	return relayLog, nil
}

// GetRelayLogs retrieves relay logs for a specific filter
func (c *Client) FilterRelayMsgs(fromBlock, toBlock *big.Int, address string) ([]*chain.RelayLog, error) {
	relayLogs := []*chain.RelayLog{}

	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{common.HexToAddress(address)},
		Topics:    RelayTopic,
	}

	rawLogs, err := c.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for _, rawLog := range rawLogs {
		relayLog, err := c.ToRelayLog(rawLog)
		if err != nil {
			return nil, err
		}
		relayLogs = append(relayLogs, relayLog)
	}

	return relayLogs, nil
}

// SubscribeToLogs subscribes to contract logs
func (c *Client) SubscribeToRelayMsgs(address string) (<-chan relay.Message, error) {
	relayMsgs := make(chan relay.Message)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(address)},
		Topics:    RelayTopic,
	}

	rawLogs := make(chan types.Log)
	sub, err := c.WsClient.SubscribeFilterLogs(context.Background(), query, rawLogs)
	if err != nil {
		c.logger.Error("Failed to subscribe to relay logs", map[string]any{
			"address": address,
			"wsurl":   c.Config.WSURL,
			"error":   err,
		})
		return nil, err
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case err := <-sub.Err():
				c.logger.Error("Relay log subscription error", map[string]any{
					"address": address,
					"wsurl":   c.Config.WSURL,
					"error":   err,
				})
				return
			case rawLog := <-rawLogs:
				relayLog, err := c.ToRelayLog(rawLog)
				if err != nil {
					c.logger.Error("Failed to parse relay log", map[string]any{
						"log":   rawLog,
						"error": err,
					})
					continue
				}
				relayMsgs <- relayLog
			}
		}
	}()

	return relayMsgs, nil
}

func (c *Client) ProcessRelayMsgs(relayMsgs <-chan relay.Message) error {
	go func() {
		select {
		case relayMsg := <-relayMsgs:
			c.logger.Info("Received new relay log", map[string]any{
				"log": relayMsg,
			})

			// TODO:
		}
	}()

	return nil
}
