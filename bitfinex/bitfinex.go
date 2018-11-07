package bitfinex

import (
	"strings"
	"sync"

	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Market contains state data
type Market struct {
	sync.RWMutex
	streaming bool
	pairs     []string
	ids       []string
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
			"tZECUSD",
			"tZECBTC",
			"tZRXUSD",
			"tZRXBTC",
		},
		ids:  []string{},
		data: make(map[string]cq.Quoter),
	}

	for _, pair := range m.pairs {
		m.ids = append(m.ids, fmtID(pair))
		m.data[pair] = Quote{
			Symbol: pair,
			ID:     fmtID(pair),
		}
	}

	return m
}

// GetIDs returns slice of pair IDs formatted with "/" (i.e., BTC/USD)
func (m *Market) GetIDs() []string {
	return m.ids
}

// Table returns an instance of tview.Table formatted for bitfinex ready for data
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

	tbl := tview.NewTable().
		SetBorders(false)

	for i, header := range headers {
		tbl.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r <= 39; r++ {
		for c := 0; c <= 8; c++ {
			tbl.SetCell(r, c, tview.NewTableCell("").
				SetAlign(tview.AlignRight))
		}
	}

	// handle errors here ***************
	m.getSnapshot()

	m.Lock()
	data := m.data
	m.Unlock()

	for _, quote := range data {
		quote.TradeUpdate(overviewTbl, tbl, tcell.AttrNone)
	}

	return tbl
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
