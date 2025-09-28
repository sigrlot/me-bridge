package chain

import (
	"sync"
	"time"

	"github.com/st-chain/me-bridge/log"
)

// Cluster manages a set of blockchain clients and provides the "best" current client.
// T must implement Client (LatestHeight). Concrete endpoints will additionally
// constrain T to relay.InEndpoint or relay.OutEndpoint via type aliases in endpoint.go.
type Cluster[T Client] struct {
	mu              sync.RWMutex
	clients         []T
	current         T
	monitorInterval time.Duration
	stopCh          chan struct{}
	logger          *log.Logger
}

// NewCluster creates a new cluster with given clients and monitor interval.
// It picks an initial current client based on the highest LatestHeight.
func NewCluster[T Client](clients []T, monitorInterval time.Duration) *Cluster[T] {
	c := &Cluster[T]{
		clients:         append([]T(nil), clients...),
		monitorInterval: monitorInterval,
		stopCh:          make(chan struct{}),
		logger:          log.WithComponent("cluster"),
	}
	c.recomputeBest()
	return c
}

// Current returns the current best client.
func (c *Cluster[T]) Current() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

// Clients returns a copy of the current client list.
func (c *Cluster[T]) Clients() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]T, len(c.clients))
	copy(out, c.clients)
	return out
}

// SetClients replaces the client set and recomputes the best current client.
func (c *Cluster[T]) SetClients(clients []T) {
	c.mu.Lock()
	c.clients = append([]T(nil), clients...)
	c.mu.Unlock()
	c.recomputeBest()
}

// ReplaceClient forces a recomputation and returns the new current client (which may be unchanged).
func (c *Cluster[T]) ReplaceClient() T {
	c.recomputeBest()
	return c.Current()
}

// Start begins background monitoring to automatically select the freshest client.
func (c *Cluster[T]) Start() {
	if c.monitorInterval <= 0 {
		return
	}
	go c.monitorLoop()
}

// Stop stops background monitoring.
func (c *Cluster[T]) Stop() {
	select {
	case <-c.stopCh:
		// already closed
	default:
		close(c.stopCh)
	}
}

// monitorLoop periodically recomputes the best client.
func (c *Cluster[T]) monitorLoop() {
	ticker := time.NewTicker(c.monitorInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.recomputeBest()
		case <-c.stopCh:
			return
		}
	}
}

// recomputeBest selects the client with the highest latest height.
func (c *Cluster[T]) recomputeBest() {
	c.mu.RLock()
	clients := append([]T(nil), c.clients...)
	c.mu.RUnlock()

	if len(clients) == 0 {
		return
	}

	var (
		best    T
		bestH   int64 = -1
		bestErr error
	)

	for _, cl := range clients {
		h, err := cl.LatestHeight()
		if err != nil {
			bestErr = err
			c.logger.Debug("latest height error", map[string]any{"error": err})
			continue
		}
		if h > bestH {
			bestH = h
			best = cl
		}
	}

	if bestH < 0 {
		// all failed; keep current but log once
		if bestErr != nil {
			c.logger.Warn("no healthy client in cluster", nil)
		}
		return
	}

	c.mu.Lock()
	c.current = best
	c.mu.Unlock()
}
