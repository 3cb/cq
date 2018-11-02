package gemini

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
	data      map[string]cq.Quoter
}

// Init creates a new instance of gemini market without quotes
func Init() *Market {
	m := &Market{
		streaming: false,
		pairs: []string{
			"btcusd",
			"ethusd",
			"ethbtc",
			// "bchusd",
			// "bchbtc",
			// "bcheth",
			"ltcusd",
			"ltcbtc",
			"zecusd",
			"zecbtc",
		},
		data: make(map[string]cq.Quoter),
	}

	for _, pair := range m.pairs {
		m.data[pair] = Quote{
			Symbol: pair,
			ID:     fmtID(pair),
		}
	}

	return m
}

// Table method uses market data to create and return an
// instance of tview.Table to caller application
func (m *Market) Table(overviewTbl *tview.Table) *tview.Table {
	headers := []string{
		"Pair",
		"Price",
		// "Change",
		"Last Size",
		"Bid",
		"Ask",
		// "Low",
		// "High",
		"Volume",
	}

	table := tview.NewTable().
		SetBorders(false)

	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r <= 17; r++ {
		for c := 0; c <= 5; c++ {
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
		quote.UpdRow(table)()
		quote.ClrBold(table)()
		quote.UpdOverviewRow(overviewTbl)()
		quote.ClrOverviewBold(overviewTbl)()
	}

	return table
}

func (m *Market) getSnapshot() []error {
	var e []error

	m.RLock()
	pairs := m.pairs
	m.RUnlock()
	l := len(pairs)
	errCh := make(chan error, l*2)

	wg := &sync.WaitGroup{}
	wg.Add(l * 2)
	for _, pair := range pairs {
		m.getTicker(pair, wg, errCh)
		m.getTrades(pair, wg, errCh)
	}
	wg.Wait()

	close(errCh)

	for err := range errCh {
		e = append(e, err)
	}

	return e
}

// Stream connects to websocket connection and starts goroutine to update state of Gemini
// market with data from websocket messages
func (m *Market) Stream(data chan cq.Quoter) error {
	err := connectWS(m, data)
	if err != nil {
		return err
	}
	return nil
}

func fmtID(pair string) string {
	temp := strings.Split(pair, "")
	temp2 := []string{temp[0], temp[1], temp[2], "/", temp[3], temp[4], temp[5]}
	return strings.ToUpper(strings.Join(temp2, ""))
}
