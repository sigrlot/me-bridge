package relay

import (
	"sync"
	"time"
)

// TransactionTracker 跟踪交易确认并处理重试
type TransactionTracker struct {
	mu                sync.RWMutex
	confirmationDepth int
	transactions      map[string]*TrackedTx // txHash -> TrackedTx
	nonceManager      *NonceManager
}

// TrackedTx 代表正在跟踪确认的交易
type TrackedTx struct {
	TxHash        string    `json:"tx_hash"`
	Nonce         uint64    `json:"nonce"`
	BlockHeight   uint64    `json:"block_height"`
	Confirmations int       `json:"confirmations"`
	CreatedAt     time.Time `json:"created_at"`
	LastChecked   time.Time `json:"last_checked"`
	Status        TxStatus  `json:"status"`
	RetryCount    int       `json:"retry_count"`
}

// NewTransactionTracker 创建一个新的交易跟踪器
func NewTransactionTracker(confirmationDepth int, nonceManager *NonceManager) *TransactionTracker {
	return &TransactionTracker{
		confirmationDepth: confirmationDepth,
		transactions:      make(map[string]*TrackedTx),
		nonceManager:      nonceManager,
	}
}

// TrackTransaction 开始跟踪一个交易
func (tt *TransactionTracker) TrackTransaction(txHash string, nonce uint64, blockHeight uint64) {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	tt.transactions[txHash] = &TrackedTx{
		TxHash:        txHash,
		Nonce:         nonce,
		BlockHeight:   blockHeight,
		Confirmations: 0,
		CreatedAt:     time.Now(),
		LastChecked:   time.Now(),
		Status:        TxStatusSubmitted,
		RetryCount:    0,
	}

	// 在nonce管理器中标记为已提交
	tt.nonceManager.MarkSubmitted(nonce, txHash)
}

// UpdateConfirmations 更新交易的确认数
func (tt *TransactionTracker) UpdateConfirmations(txHash string, currentBlockHeight uint64) bool {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	tx, exists := tt.transactions[txHash]
	if !exists {
		return false
	}

	if currentBlockHeight > tx.BlockHeight {
		tx.Confirmations = int(currentBlockHeight - tx.BlockHeight)
		tx.LastChecked = time.Now()

		// 检查交易是否已确认
		if tx.Confirmations >= tt.confirmationDepth {
			tx.Status = TxStatusConfirmed
			// 在nonce管理器中标记为已确认并从跟踪中移除
			tt.nonceManager.MarkConfirmed(tx.Nonce)
			delete(tt.transactions, txHash)
			return true
		}
	}

	return false
}

// MarkFailed 将交易标记为失败
func (tt *TransactionTracker) MarkFailed(txHash string, reason string) {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	tx, exists := tt.transactions[txHash]
	if !exists {
		return
	}

	tx.Status = TxStatusFailed
	tx.RetryCount++
	tt.nonceManager.MarkFailed(tx.Nonce)
}

// GetPendingTransactions 返回所有待处理交易
func (tt *TransactionTracker) GetPendingTransactions() []*TrackedTx {
	tt.mu.RLock()
	defer tt.mu.RUnlock()

	var pending []*TrackedTx
	for _, tx := range tt.transactions {
		if tx.Status == TxStatusSubmitted {
			txCopy := *tx
			pending = append(pending, &txCopy)
		}
	}

	return pending
}

// CleanupStale 移除最近没有更新的过期交易
func (tt *TransactionTracker) CleanupStale(maxAge time.Duration) int {
	tt.mu.Lock()
	defer tt.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	cleaned := 0

	for txHash, tx := range tt.transactions {
		if tx.LastChecked.Before(cutoff) {
			delete(tt.transactions, txHash)
			cleaned++
		}
	}

	return cleaned
}

// GetTransactionStatus 返回特定交易的状态
func (tt *TransactionTracker) GetTransactionStatus(txHash string) (*TrackedTx, bool) {
	tt.mu.RLock()
	defer tt.mu.RUnlock()

	tx, exists := tt.transactions[txHash]
	if !exists {
		return nil, false
	}

	txCopy := *tx
	return &txCopy, true
}
