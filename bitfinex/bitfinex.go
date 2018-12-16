package bitfinex

import (
	"strings"
	"sync"

	"github.com/3cb/cq/cq"
)

// Market contains state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	ids       []string
	data      map[string]cq.Quote
}

// Init creates instance of a bitfinex market without quotes
func Init() *Market {
	m := &Market{
		streaming: false,
		pairs: []string{
			"tBTCUSD",
			"tBTCEUR",
			"tBTCGBP",
			"tBTCJPY",
			"tBCHUSD",
			"tBCHBTC",
			"tETHUSD",
			"tETHBTC",
			"tETHEUR",
			"tETHGBP",
			"tETHJPY",
			"tLTCUSD",
			"tLTCBTC",
			"tZECUSD",
			"tZECBTC",
			"tZRXUSD",
			"tZRXBTC",
		},
		ids:  []string{},
		data: make(map[string]cq.Quote),
	}

	for _, pair := range m.pairs {
		m.ids = append(m.ids, fmtID(pair))
		m.data[pair] = cq.Quote{
			MarketID: "bitfinex",
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

// getSnapshot performs http requests to the Bitfinex API to get initial market data
func (m *Market) getSnapshot() []error {
	var e []error

	m.RLock()
	pairs := m.pairs
	m.RUnlock()
	l := len(pairs)
	errCh := make(chan error, l)

	wg := &sync.WaitGroup{}
	wg.Add(l + 1)
	go m.getTickers(errCh, wg)
	for _, pair := range pairs {
		go m.getTrades(pair, errCh, wg)
	}
	wg.Wait()

	close(errCh)
	for err := range errCh {
		e = append(e, err)
	}
	return e
}

// Stream connects to Bitfinex websocket API and sends messages to data channel
// to update market and overview tables
func (m *Market) Stream(timerCh chan<- cq.TimerMsg) error {
	err := connectWS(m, timerCh)
	if err != nil {
		return err
	}
	return nil
}

func fmtID(symbol string) string {
	s1 := strings.Split(symbol, "")
	b := strings.Builder{}
	for i := 1; i < 4; i++ {
		b.WriteString(s1[i])
	}
	b.WriteString("/")
	for i := 4; i < 7; i++ {
		b.WriteString(s1[i])
	}
	return b.String()
}
