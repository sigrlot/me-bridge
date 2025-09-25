package relay


// FeeCalculator 计算跨链交易的费用
type FeeCalculator interface {
	CalculateFee(amount int64) int64
}

