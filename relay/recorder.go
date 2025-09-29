package relay

import (
	"sync"
	"time"

	logger "github.com/st-chain/me-bridge/log"
)

// TxRecorder 管理异步交易的nonce分配
type TxRecorder struct {
	mu           sync.RWMutex
	currentNonce uint64
	pendingTxs   map[uint64]*PendingTx // nonce -> pending transaction
	logger       *logger.Logger
}

// PendingTx 代表一个帶有nonce的待处理交易
type PendingTx struct {
	Nonce     uint64    `json:"nonce"`
	TxHash    string    `json:"tx_hash,omitempty"` // 提交后填充
	Msg       *OutMsg   `json:"msg"`
	CreatedAt time.Time `json:"created_at"`
	Retries   int       `json:"retries"`
	Status    TxStatus  `json:"status"`
}

// TxStatus 代表交易状态
type TxStatus int

const (
	TxStatusPending TxStatus = iota
	TxStatusSubmitted
	TxStatusConfirmed
	TxStatusFailed
)

func (s TxStatus) String() string {
	switch s {
	case TxStatusPending:
		return "pending"
	case TxStatusSubmitted:
		return "submitted"
	case TxStatusConfirmed:
		return "confirmed"
	case TxStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// NewTxRecorder 创建一个从指定nonce开始的新nonce管理器
func NewTxRecorder(startNonce uint64) *TxRecorder {
	return &TxRecorder{
		currentNonce: startNonce,
		pendingTxs:   make(map[uint64]*PendingTx),
		logger:       logger.WithComponent("nonce-manager"),
	}
}

// AllocateNonce 为消息分配新的nonce并进行跟踪
func (nm *TxRecorder) AllocateNonce(msg *OutMsg) uint64 {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	nonce := nm.currentNonce
	nm.currentNonce++

	pendingTx := &PendingTx{
		Nonce:     nonce,
		Msg:       msg,
		CreatedAt: time.Now(),
		Status:    TxStatusPending,
	}

	nm.pendingTxs[nonce] = pendingTx

	nm.logger.Debug("allocated nonce", map[string]any{
		"nonce":    nonce,
		"sender":   msg.Sender,
		"receiver": msg.Receiver,
		"amount":   msg.Amount,
	})

	return nonce
}

// MarkSubmitted 将交易标记为已提交并记录其哈希
func (nm *TxRecorder) MarkSubmitted(nonce uint64, txHash string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	tx, exists := nm.pendingTxs[nonce]
	if !exists {
		nm.logger.Warn("尝试将不存在的nonce标记为已提交", map[string]any{
			"nonce":   nonce,
			"tx_hash": txHash,
		})
		return ErrNonceNotFound
	}

	tx.TxHash = txHash
	tx.Status = TxStatusSubmitted

	nm.logger.Debug("将交易标记为已提交", map[string]any{
		"nonce":   nonce,
		"tx_hash": txHash,
	})

	return nil
}

// MarkConfirmed 将交易标记为已确认并从待处理列表中移除
func (nm *TxRecorder) MarkConfirmed(nonce uint64) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	tx, exists := nm.pendingTxs[nonce]
	if !exists {
		nm.logger.Warn("尝试确认不存在的nonce", map[string]any{
			"nonce": nonce,
		})
		return ErrNonceNotFound
	}

	tx.Status = TxStatusConfirmed
	delete(nm.pendingTxs, nonce)

	nm.logger.Debug("已确认并移除交易", map[string]any{
		"nonce":   nonce,
		"tx_hash": tx.TxHash,
	})

	return nil
}

// MarkFailed 将交易标记为失败
func (nm *TxRecorder) MarkFailed(nonce uint64) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	tx, exists := nm.pendingTxs[nonce]
	if !exists {
		return ErrNonceNotFound
	}

	tx.Status = TxStatusFailed
	tx.Retries++

	nm.logger.Warn("将交易标记为失败", map[string]any{
		"nonce":   nonce,
		"tx_hash": tx.TxHash,
		"retries": tx.Retries,
	})

	return nil
}

// GetPendingTx 根据nonce返回待处理交易
func (nm *TxRecorder) GetPendingTx(nonce uint64) (*PendingTx, bool) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	tx, exists := nm.pendingTxs[nonce]
	if !exists {
		return nil, false
	}

	// 返回副本以避免竞态条件
	txCopy := *tx
	return &txCopy, true
}

// GetCurrentNonce 返回当前nonce（下一个将被分配的）
func (nm *TxRecorder) GetCurrentNonce() uint64 {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return nm.currentNonce
}

// GetPendingCount 返回待处理交易的数量
func (nm *TxRecorder) GetPendingCount() int {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return len(nm.pendingTxs)
}

// CleanupStale 移除超过给定时间的过期待处理交易
func (nm *TxRecorder) CleanupStale(maxAge time.Duration) int {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	cleaned := 0

	for nonce, tx := range nm.pendingTxs {
		if tx.CreatedAt.Before(cutoff) && tx.Status != TxStatusSubmitted {
			delete(nm.pendingTxs, nonce)
			cleaned++
			nm.logger.Debug("清理过期交易", map[string]any{
				"nonce":      nonce,
				"created_at": tx.CreatedAt,
				"status":     tx.Status.String(),
			})
		}
	}

	if cleaned > 0 {
		nm.logger.Info("清理过期交易", map[string]any{
			"cleaned_count": cleaned,
			"remaining":     len(nm.pendingTxs),
		})
	}

	return cleaned
}

// GetRetryableTxs 返回可以重试的交易
func (nm *TxRecorder) GetRetryableTxs(maxRetries int) []*PendingTx {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	var retryable []*PendingTx
	for _, tx := range nm.pendingTxs {
		if tx.Status == TxStatusFailed && tx.Retries < maxRetries {
			txCopy := *tx
			retryable = append(retryable, &txCopy)
		}
	}

	return retryable
}
