package bsc

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

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

// ProcessInMsgs 处理跨链消息，返回错误通道供 Tunnel 监听
func (c *Client) ProcessInMsgs(msgs <-chan relay.InMsg) <-chan relay.ProcessError {
	errorChan := make(chan relay.ProcessError, 10)

	go func() {
		defer close(errorChan)

		for msg := range msgs {
			if err := c.processMessage(msg); err != nil {
				// 判断错误类型
				recoverable := relay.IsRecoverable(err)

				if recoverable {
					// Target 自行重试
					if retryErr := c.retryMessage(msg, 3); retryErr != nil {
						// 重试失败，上报给 Tunnel
						errorChan <- relay.ProcessError{
							Msg:         msg,
							Err:         retryErr,
							Recoverable: false,
							Timestamp:   time.Now(),
							RetryCount:  3,
						}
					}
				} else {
					// 直接上报严重错误
					errorChan <- relay.ProcessError{
						Msg:         msg,
						Err:         err,
						Recoverable: false,
						Timestamp:   time.Now(),
						RetryCount:  0,
					}
				}
			}
		}
	}()

	return errorChan
}

// processMessage 处理单个跨链消息
func (c *Client) processMessage(msg relay.InMsg) error {
	c.logger.Info("Processing cross-chain message", map[string]any{
		"msg": msg,
	})

	// TODO: 实现具体的交易发送逻辑
	// 1. 构造交易
	// 2. 签名交易
	// 3. 发送交易
	// 4. 处理发送结果

	return nil // 临时返回，实际实现时需要处理具体逻辑
}

// retryMessage 重试消息处理
func (c *Client) retryMessage(msg relay.InMsg, maxRetries int) error {
	strategy := relay.NewDefaultRecoveryStrategy()

	for i := 0; i < maxRetries; i++ {
		if err := c.processMessage(msg); err == nil {
			return nil // 成功
		}

		// 使用恢复策略等待
		procErr := relay.ProcessError{
			Msg:        msg,
			Err:        relay.ErrProcessingFailed,
			RetryCount: i,
		}

		if retryErr := strategy.Recover(procErr); retryErr != nil {
			return retryErr
		}
	}

	return relay.ErrProcessingFailed
}

// Reset 重置客户端状态，用于错误恢复
func (c *Client) Reset() error {
	c.logger.Info("Resetting BSC client state")

	// 重新连接客户端
	if c.Client != nil {
		c.Client.Close()
	}
	if c.WsClient != nil {
		c.WsClient.Close()
	}

	// 重新建立连接
	client, err := ethclient.Dial(c.Config.RPCURL)
	if err != nil {
		return err
	}

	wsClient, err := ethclient.Dial(c.Config.WSURL)
	if err != nil {
		client.Close()
		return err
	}

	c.Client = client
	c.WsClient = wsClient

	// 重新同步状态
	height, err := c.GetLatestHeight()
	if err != nil {
		return err
	}
	c.latestHeight = height

	return nil
}

// GetSequence 获取当前序列号和高度
func (c *Client) GetSequence() (uint64, uint64) {
	// TODO: 从合约或数据库获取实际的序列号
	// 这里返回模拟数据
	return 0, c.latestHeight
}

// GetCurrentNonce 获取当前用于跨链的 nonce (实现 OutEndpoint 接口)
func (c *Client) GetCurrentNonce() uint64 {
	// TODO: 从实际的 nonce 管理器获取
	// 这里返回模拟数据
	return 0
}
