package relay

// Checker 检查跨链桥状态，并发出预警信息
type Checker interface {
	BalanceWarning() error
}
