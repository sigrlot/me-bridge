package bsc

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// SendTransaction signs and sends a transaction to BSC
func (c *Client) SendTransaction(tx *Transaction) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Create the transaction
	ethTx := types.NewTransaction(
		tx.Nonce,
		tx.To,
		tx.Value,
		tx.GasLimit,
		tx.GasPrice,
		tx.Data,
	)

	// Sign the transaction
	signedTx, err := types.SignTx(ethTx, types.NewEIP155Signer(c.chainID), c.privateKey)
	if err != nil {
		return nil, err
	}

	// Send the transaction
	err = c.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// CallContract calls a smart contract method (read-only)
func (c *Client) CallContract(contractAddress common.Address, data []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	return c.Client.CallContract(ctx, msg, nil)
}

// WaitForTransactionReceipt waits for transaction to be mined
func (c *Client) WaitForTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	for {
		receipt, err := c.Client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
			// Continue polling
		}
	}
}
