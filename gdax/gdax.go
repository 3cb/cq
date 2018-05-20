package gdax

import (
	"fmt"
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

	for r := 1; r < len(m.pairs)+1; r++ {
		pair := m.pairs[r-1]
		quote := m.data[pair]
		delta, color := display.FmtDelta(quote.Price, quote.Open)

		table.SetCell(r, 0, tview.NewTableCell(display.FmtPair(pair)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 1, tview.NewTableCell(display.FmtPrice(quote.Price)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 2, tview.NewTableCell(delta).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 3, tview.NewTableCell(display.FmtSize(quote.Size)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 4, tview.NewTableCell(display.FmtPrice(quote.Bid)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 5, tview.NewTableCell(display.FmtPrice(quote.Ask)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 6, tview.NewTableCell(display.FmtPrice(quote.Low)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 7, tview.NewTableCell(display.FmtPrice(quote.High)).
			SetTextColor(color).
			SetAlign(tview.AlignRight))
		table.SetCell(r, 8, tview.NewTableCell(display.FmtVolume(quote.Volume)).
			SetTextColor(color).
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
