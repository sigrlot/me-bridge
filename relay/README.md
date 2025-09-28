# Relay Package Implementation

This package implements a cross-chain bridge relay system with comprehensive nonce management for async transactions.

## Key Components

### 1. Relay (`relay.go`)

The main relay orchestrator that connects source and target endpoints:

- **Message Routing**: Routes inbound messages from source to target chain
- **Async Processing**: Handles message processing asynchronously with proper goroutine management
- **Nonce Management**: Integrates with NonceManager for transaction sequencing
- **Fee Calculation**: Supports pluggable fee calculation strategies
- **Graceful Shutdown**: Supports clean shutdown with `Stop()` method

### 2. Nonce Manager (`nonce.go`)

Critical for async transaction handling:

- **Nonce Allocation**: Thread-safe nonce assignment for outbound transactions
- **Transaction Tracking**: Tracks pending transactions by nonce
- **State Management**: Handles pending/submitted/confirmed/failed states
- **Retry Logic**: Supports transaction retry with configurable limits
- **Cleanup**: Automatic cleanup of stale transactions

### 3. Transaction Tracker (`tracker.go`)

Monitors transaction confirmations:

- **Confirmation Tracking**: Tracks block confirmations for transactions
- **Status Updates**: Updates transaction status based on chain state
- **Failure Handling**: Marks failed transactions and coordinates with nonce manager
- **Stale Cleanup**: Removes old unconfirmed transactions

### 4. Interfaces (`endpoint.go`)

Defines contracts for chain endpoints:

- **InEndpoint**: Handles inbound messages (subscribe + process outbound confirmations)
- **OutEndpoint**: Handles outbound messages (process inbound + subscribe to confirmations)
- **Consistent API**: Standardized channel-based communication

## Key Features for Async Transactions

### Nonce Management

```go
// Allocate nonce for new transaction
nonce := relay.NonceManager.AllocateNonce(outMsg)

// Mark as submitted when transaction is sent
relay.NonceManager.MarkSubmitted(nonce, txHash)

// Mark as confirmed when enough confirmations received
relay.NonceManager.MarkConfirmed(nonce)
```

### Transaction Flow

1. **Inbound Message**: Source chain message received via subscription
2. **Nonce Allocation**: Assign sequential nonce for target chain transaction
3. **Async Submission**: Submit transaction to target chain asynchronously
4. **Confirmation Tracking**: Monitor transaction confirmations
5. **State Update**: Update nonce manager and clean up on confirmation

### Error Handling & Retries

- Failed transactions are marked in nonce manager
- Configurable retry limits with exponential backoff
- Stale transaction cleanup prevents nonce gaps
- Graceful degradation on endpoint failures

## Usage Example

```go
// Initialize components
nonceManager := NewNonceManager(startNonce)
feeCalculator := &MyFeeCalculator{}
config := &RelayConfig{...}

// Create relay
relay := NewRelay(config, sourceEndpoint, targetEndpoint, feeCalculator, startNonce)

// Start message processing
if err := relay.Work(); err != nil {
    log.Fatal(err)
}

// Monitor status
status := relay.GetStatus()
fmt.Printf("Pending transactions: %d\n", status["pending_count"])
```

## Design Principles

1. **Async-First**: All transaction operations are non-blocking
2. **Nonce Safety**: Strict nonce sequencing prevents transaction conflicts
3. **Observability**: Rich status reporting and logging
4. **Fault Tolerance**: Graceful handling of chain downtime and failures
5. **Resource Management**: Automatic cleanup of stale data

This implementation ensures reliable cross-chain message relay even under adverse conditions like network partitions, chain reorganizations, and temporary endpoint failures.
