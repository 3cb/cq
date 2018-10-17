package bitfinex

import (
	"sync"

	"github.com/3cb/cq/cq"
	"github.com/3cb/muttview"
	"github.com/gdamore/tcell"
)

// Market contains state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	sockets   map[string]int
	data      map[string]cq.Quoter
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
		},
		sockets: make(map[string]int),
		data:    make(map[string]cq.Quoter),
	}

	for _, pair := range m.pairs {
		m.data[pair] = Quote{}
	}

	return m
}

// GetSnapshot performs http requests to the Bitfinex API to get initial market data
func (m *Market) GetSnapshot() []error {
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

// Table returns an instance of tview.Table formatted for bitfinex ready for data
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

	tbl := tview.NewTable().
		SetBorders(false)

	for i, header := range headers {
		tbl.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r <= 33; r++ {
		for c := 0; c <= 8; c++ {
			tbl.SetCell(r, c, tview.NewTableCell("").
				SetAlign(tview.AlignRight))
		}
	}

	return tbl
}

// PrimeTables ranges over data map of price Quotes and sends to data channel
func (m *Market) PrimeTables(data chan cq.Quoter) {
	m.RLock()
	quotes := m.data
	m.RUnlock()
	for _, v := range quotes {
		data <- v
	}
}

// Stream connects to Bitfinex websocket API and sends messages to data channel
// to update market and overview tables
func (m *Market) Stream(data chan cq.Quoter) error {
	err := connectWS(m, data)
	if err != nil {
		return err
	}
	return nil
}
