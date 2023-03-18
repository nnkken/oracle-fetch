package types

import (
	"math/big"
	"time"
)

type DBEntry struct {
	Token string
	Unit  string
	// Price indicate the price of token in 1e-8 unit, i.e.
	// Token = ETH, Unit = USD, Price = 160000000000 means 1ETH = 1600 USD
	Price          *big.Int
	PriceTimestamp time.Time
	FetchTimestamp time.Time
	// TODO: metadata (JSON?)
}

type DataSource interface {
	Fetch() ([]DBEntry, error)
}
