package gdax

import (
	"fmt"
	"sync"

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

// Quote contains most recent data for each crypto currency pair
// Trailing comments denote which http request or websocket stream the data comes from
// getTrades: https://docs.gdax.com/#get-trades
// match: https://docs.gdax.com/#the-code-classprettyprintfullcode-channel
// ticker: https://docs.gdax.com/#the-code-classprettyprinttickercode-channel
// *** GDAX API documentation for websocket ticker channel does not show all available fields as of 2/11/2018
type Quote struct {
	Type string `json:"type"` // used to filter websocket messages

	ID    string `json:"product_id"`
	Price string `json:"price"` // getTrades/match
	Size  string `json:"size"`  // getTrades/match

	Delta string // % change in price

	Bid string `json:"best_bid"` // getTicker/ticker
	Ask string `json:"best_ask"` // getTicker/ticker

	High   string `json:"high_24h"`   // getStats/ticker
	Low    string `json:"low_24h"`    // getStats/ticker
	Open   string `json:"open_24h"`   // getStats/ticker
	Volume string `json:"volume_24h"` // getStats/ticker
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

func (m *Market) Stream() error {
	err := connectWS(m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Market) Display() {

}

func (m *Market) Table() *tview.Table {
	headers := []string{
		"Pair",
		"Price",
		"Change",
		"Last Size",
		"Bid",
		"Ask",
		"High",
		"Low",
		"Volume",
	}

	table := tview.NewTable().
		SetBorders(false)

	// col, rows := 9, 12
	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r < len(m.pairs)+1; r++ {
		pair := m.pairs[r-1]
		table.SetCell(r, 0, tview.NewTableCell(pair).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 1, tview.NewTableCell(m.data[pair].Price).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 2, tview.NewTableCell(m.data[pair].Delta).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 3, tview.NewTableCell(m.data[pair].Size).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 4, tview.NewTableCell(m.data[pair].Bid).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 5, tview.NewTableCell(m.data[pair].Ask).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 6, tview.NewTableCell(m.data[pair].High).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 7, tview.NewTableCell(m.data[pair].Low).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 8, tview.NewTableCell(m.data[pair].Volume).
			// SetTextColor(headerColor).
			SetAlign(tview.AlignRight))
	}
	return table
}

func (m *Market) Print() {
	m.RLock()
	for k, v := range m.data {
		fmt.Printf("%v: %+v", k, v)
	}
	m.RUnlock()
}
