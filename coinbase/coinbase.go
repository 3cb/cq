package coinbase

import (
	"strings"
	"sync"

	"github.com/3cb/cq/cq"
)

// Market conatins state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	ids       []string
	data      map[string]cq.Quote
}

// Init initializes and returns an instance of the Coinbase exchange without quotes
func Init() *Market {
	m := &Market{
		streaming: false,
		pairs: []string{
			"BTC-USD",
			"BTC-EUR",
			"BTC-GBP",
			"BCH-USD",
			"BCH-BTC",
			"BCH-EUR",
			"ETH-USD",
			"ETH-BTC",
			"ETH-EUR",
			"ETH-GBP",
			"LTC-USD",
			"LTC-BTC",
			"LTC-EUR",
			"ZRX-USD",
			"ZRX-BTC",
		},
		ids:  []string{},
		data: make(map[string]cq.Quote),
	}

	// Create slice of ids formatted based on pair strings
	// Instantiate quote for each pair with MarketID and PairID already set
	for _, pair := range m.pairs {
		m.ids = append(m.ids, fmtID(pair))
		m.data[pair] = cq.Quote{
			MarketID: "coinbase",
			ID:       fmtID(pair),
		}
	}

	return m
}

// GetIDs returns slice of pair IDs formatted with "/" (i.e., BTC/USD)
func (m *Market) GetIDs() []string {
	return m.ids
}

// GetQuotes returns a map used to prime overview table with data
// Keys are pair IDs separated with "/". Values are of type Quote.
func (m *Market) GetQuotes() map[string]cq.Quote {
	m.Lock()
	d := m.data
	m.Unlock()
	return d
}

// Stream connects to websocket connection and starts goroutine to update state of GDAX
// market with data from websocket messages
func (m *Market) Stream(routerCh chan<- cq.TimerMsg) error {
	err := connectWS(m, routerCh)
	if err != nil {
		return err
	}
	return nil
}

// fmtID formats product id to represent currency pair (i.e., "BTC/USD")
func fmtID(id string) string {
	return strings.Join(strings.Split(id, "-"), "/")
}
