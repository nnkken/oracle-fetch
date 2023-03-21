package utils

import (
	"math/big"
)

const Decimals = 8

// ComputeDecimalShift compute the decimal shift between the fetched price and the price we want to store
func ComputeDecimalShift(fetchedDecimals uint8) int {
	return Decimals - int(fetchedDecimals)
}

// NormalizePrice utilize ComputeDecimalShift to shift the price
// i.e. storedPrice = fetchedPrice * 10^decimalShift
func NormalizePrice(price *big.Int, fetchedDecimals uint8) *big.Int {
	decimalShift := ComputeDecimalShift(fetchedDecimals)
	if decimalShift == 0 {
		return price
	}
	if decimalShift > 0 {
		return new(big.Int).Mul(price, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalShift)), nil))
	}
	// TODO: maybe not integer? do we care?
	return new(big.Int).Div(price, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-decimalShift)), nil))
}
