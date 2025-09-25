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

// Transaction represents a BSC transaction
type Transaction struct {
	To       common.Address
	Value    *big.Int
	Data     []byte
	GasLimit uint64
	GasPrice *big.Int
	Nonce    uint64
}

// EstimateTransactionGas estimates gas for a BSC transaction
func (c *BSCClient) EstimateTransactionGas(tx *Transaction) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	fromAddress := crypto.PubkeyToAddress(c.privateKey.PublicKey)

	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &tx.To,
		Value:    tx.Value,
		Data:     tx.Data,
		GasPrice: tx.GasPrice,
	}

	return c.client.EstimateGas(ctx, msg)
}

// SendTransaction signs and sends a transaction to BSC
func (c *BSCClient) SendTransaction(tx *Transaction) (*types.Transaction, error) {
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
	err = c.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// SendBNB sends BNB to an address
func (c *BSCClient) SendBNB(to common.Address, amount *big.Int) (*types.Transaction, error) {
	fromAddress := crypto.PubkeyToAddress(c.privateKey.PublicKey)

	nonce, err := c.GetNonce(fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.GetGasPrice()
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		To:       to,
		Value:    amount,
		Data:     nil,
		GasLimit: 21000, // Standard gas limit for BNB transfer
		GasPrice: gasPrice,
		Nonce:    nonce,
	}

	return c.SendTransaction(tx)
}

// CallContract calls a smart contract method (read-only)
func (c *BSCClient) CallContract(contractAddress common.Address, data []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	return c.client.CallContract(ctx, msg, nil)
}

// WaitForTransactionReceipt waits for transaction to be mined
func (c *BSCClient) WaitForTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	for {
		receipt, err := c.client.TransactionReceipt(ctx, txHash)
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
