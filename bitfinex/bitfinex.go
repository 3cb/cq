package bitfinex

import (
	"sync"

	"github.com/3cb/cq/cq"
)

// Market contains state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	data      map[string]cq.Quoter
}

// Init creates instance of a bitfinex market without quotes
func Init() *Market {
	m := &Market{
		streaming: false,
		pairs: []string{
			"BTC-USD",
			"BTC-EUR",
			"BTC-GBP",
			"BTC-JPY",
			"BCH-USD",
			"BCH-BTC",
			"BCH-EUR",
			"ETH-USD",
			"ETH-BTC",
			"ETH-EUR",
			"ETH-GBP",
			"ETH-JPY",
			"LTC-USD",
			"LTC-BTC",
			"LTC-EUR",
		},
		data: make(map[string]cq.Quoter),
	}

	// for _, pair := range m.pairs {
	// 	m.data[pair] = Quote{}
	// }

	return m
}
