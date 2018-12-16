package hitbtc

import (
	"strings"
	"sync"

	"github.com/3cb/cq/cq"
)

// Market contains state data
type Market struct {
	sync.RWMutex
	streaming bool
	symbols   []string
	ids       []string
	data      map[string]cq.Quote
}

// Init creates instance of a bitfinex market without quotes
func Init() *Market {
	m := &Market{
		streaming: false,
		symbols: []string{
			"BTCUSD",
			"BCHUSD",
			"ETHUSD",
			"ETHBTC",
			"LTCUSD",
			"LTCBTC",
			"ZECUSD",
			"ZECBTC",
			"ZRXUSD",
			"ZRXBTC",
		},
		ids:  []string{},
		data: make(map[string]cq.Quote),
	}

	for _, s := range m.symbols {
		m.ids = append(m.ids, fmtID(s))
		m.data[s] = cq.Quote{
			MarketID: "hitbtc",
			ID:       fmtID(s),
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

// Stream connects to websocket server and streams price quotes
func (m *Market) Stream(timerCh chan<- cq.TimerMsg) error {
	if err := connectWS(m, timerCh); err != nil {
		return err
	}

	return nil
}

// formats pair with uppercase letters separated by "/"
func fmtID(s string) string {
	temp := strings.Split(s, "")
	b := strings.Builder{}

	for i := 0; i < 3; i++ {
		b.WriteString(temp[i])
	}
	b.WriteString("/")
	for i := 3; i < len(temp); i++ {
		b.WriteString(temp[i])
	}

	return b.String()
}
