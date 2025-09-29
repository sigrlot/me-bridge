package relay

import "math/big"

// FeeCalculator 计算跨链交易的费用
type FeeCalculator struct{}

func (fc *FeeCalculator) CalculateFee(value *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}
