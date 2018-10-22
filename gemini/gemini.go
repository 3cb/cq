package gemini

import (
	"sync"

	"github.com/3cb/cq/cq"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
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
		pair: []string{
			"btcusd",
			"ethusd",
			"ethbtc",
			"bchusd",
			"bchbtc",
			"ltcusd",
			"ltcbtc",
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
	// errCh := make(chan error, ())

}
