package coinbase

import (
	"strings"
	"sync"
	"time"

	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Market conatins state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	ids       []string
	data      map[string]cq.Quoter
}

// Init initializes and returns an instance of the GDAX exchange without quotes
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
		data: make(map[string]cq.Quoter),
	}

	for _, pair := range m.pairs {
		m.ids = append(m.ids, fmtID(pair))
		m.data[pair] = Quote{
			ID: fmtID(pair),
		}
	}

	return m
}

// GetIDs returns slice of pair IDs formatted with "/" (i.e., BTC/USD)
func (m *Market) GetIDs() []string {
	return m.ids
}

// Table method uses market data to create and return an
// instance of tview.Table to caller application
func (m *Market) Table(overviewTbl *tview.Table) *tview.Table {
	headers := []string{
		"Pair",
		"Price",
		"Change",
		"Last Size",
		"Bid",
		"Ask",
		"Low",
		"High",
		"Volume",
	}

	table := tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorLightSlateGray)

	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r < len(m.pairs)+1; r++ {
		for c := 0; c <= 8; c++ {
			table.SetCell(r, c, tview.NewTableCell("-").
				SetAlign(tview.AlignRight))
		}
	}

	// handle errors here ***************
	m.getSnapshot()

	m.Lock()
	data := m.data
	m.Unlock()

	for _, quote := range data {
		quote.InsertTrade(overviewTbl, table, tcell.AttrNone)()
	}

	return table
}

// getSnapshot method performs concurrent http requests the GDAX REST API to get initial
// market data
func (m *Market) getSnapshot() []error {
	var e []error

	m.RLock()
	pairs := m.pairs
	m.RUnlock()
	l := len(pairs)
	errCh := make(chan error, (9 * l))

	wg := &sync.WaitGroup{}
	wg.Add(3 * l)
	// break requests up into bursts to satisfy coinbase throttling
	for i := 0; i < 5; i++ {
		go getTrades(m, pairs[i], wg, errCh)
		go getStats(m, pairs[i], wg, errCh)
		go getTicker(m, pairs[i], wg, errCh)
	}
	time.Sleep(1 * time.Second)
	for i := 5; i < 10; i++ {
		go getTrades(m, pairs[i], wg, errCh)
		go getStats(m, pairs[i], wg, errCh)
		go getTicker(m, pairs[i], wg, errCh)
	}
	time.Sleep(1 * time.Second)
	for i := 10; i < 15; i++ {
		go getTrades(m, pairs[i], wg, errCh)
		go getStats(m, pairs[i], wg, errCh)
		go getTicker(m, pairs[i], wg, errCh)
	}
	wg.Wait()

	close(errCh)
	for err := range errCh {
		e = append(e, err)
	}

	return e
}

// Stream connects to websocket connection and starts goroutine to update state of GDAX
// market with data from websocket messages
func (m *Market) Stream(timerCh chan<- cq.TimerMsg) error {
	err := connectWS(m, timerCh)
	if err != nil {
		return err
	}
	return nil
}

// fmtID formats product id to represent currency pair (i.e., "BTC/USD")
func fmtID(id string) string {
	return strings.Join(strings.Split(id, "-"), "/")
}
