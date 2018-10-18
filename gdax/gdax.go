package gdax

import (
	"sync"

	"github.com/3cb/cq/cq"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
)

// Market conatins state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
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
			"LTC-USD",
			"LTC-BTC",
			"LTC-EUR",
		},
		data: make(map[string]cq.Quoter),
	}

	for _, pair := range m.pairs {
		m.data[pair] = Quote{}
	}

	return m
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
		SetBorders(false)

	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r <= 27; r++ {
		for c := 0; c <= 8; c++ {
			table.SetCell(r, c, tview.NewTableCell("").
				SetAlign(tview.AlignRight))
		}
	}

	// handle errors here ***************
	m.getSnapshot()

	m.Lock()
	data := m.data
	m.Unlock()

	for _, quote := range data {
		quote.UpdRow(table)
		quote.ClrBold(table)
		quote.UpdOverviewRow(overviewTbl)
		quote.ClrOverviewBold(overviewTbl)
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
	for _, pair := range pairs {
		go getTrades(m, pair, wg, errCh)
		go getStats(m, pair, wg, errCh)
		go getTicker(m, pair, wg, errCh)
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
func (m *Market) Stream(data chan cq.Quoter) error {
	err := connectWS(m, data)
	if err != nil {
		return err
	}
	return nil
}
