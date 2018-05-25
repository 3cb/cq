package gdax

import (
	"sync"

	"github.com/3cb/cq/display"
	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

// Market conatins state data for the GDAX market
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	data      map[string]Quote
}

// Init initializes and returns an instance of the GDAX exchange
func Init() *Market {
	return &Market{
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
		data: make(map[string]Quote),
	}
}

// Snapshot method performs concurrent http requests the GDAX REST API to get initial
// market data
func (m *Market) Snapshot() []error {
	var e []error
	errCh := make(chan error, (9 * len(m.pairs)))

	wg := &sync.WaitGroup{}
	wg.Add(3 * len(m.pairs))
	for _, pair := range m.pairs {
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
func (m *Market) Stream(upd chan display.Setter) error {
	err := connectWS(m, upd)
	if err != nil {
		return err
	}
	return nil
}

// func (m *Market) Display() {

// }

// Table method uses maket data to create and return an
// instance of tview.Table to caller application
func (m *Market) Table() *tview.Table {
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

	for r := 0; r < len(m.pairs); r++ {
		m.RLock()
		pair := m.pairs[r]
		quote := m.data[pair]
		m.RUnlock()
		quote.SetRow(table)
	}
	return table
}

// func (m *Market) Print() {
// 	m.RLock()
// 	for k, v := range m.data {
// 		fmt.Printf("%v: %+v", k, v)
// 	}
// 	m.RUnlock()
// }
