package bsc

import "errors"

// nonce已经被使用
var ErrNonceUsed = errors.New("nonce already used")

// 余额不足
var ErrInsufficientBalance = errors.New("insufficient balance")

// 交易未找到
var ErrTransactionNotFound = errors.New("transaction not found")

// 交易失败
var ErrTransactionFailed = errors.New("transaction failed")

// 节点不可用
var ErrNodeUnavailable = errors.New("node unavailable")

// 超时错误
var ErrTimeout = errors.New("operation timed out")

// 未知错误
var ErrUnknown = errors.New("unknown error")

// 无效地址
var ErrInvalidAddress = errors.New("invalid address")

// 无效交易
var ErrInvalidTransaction = errors.New("invalid transaction")
