package hitbtc

import (
	"strings"
	"sync"

	"github.com/3cb/cq/cq"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
)

// Market contains state data
type Market struct {
	sync.RWMutex
	streaming bool
	symbols   []string
	data      map[string]cq.Quoter
}

// Init creates instance of a bitfinex market without quotes
func Init() *Market {
	m := &Market{
		streaming: false,
		symbols: []string{
			"BTCUSD",
			"BCHUSD",
			"ETHUSD",
			"ETHBTC",
			"LTCUSD",
			"LTCBTC",
			"ZECUSD",
			"ZECBTC",
			"ZRXUSD",
			"ZRXBTC",
		},
		data: make(map[string]cq.Quoter),
	}

	for _, s := range m.symbols {
		m.data[s] = Quote{
			Symbol: s,
			ID:     fmtID(s),
		}
	}

	return m
}

// Table returns an instance of tview.Table formatted for hitbtc ready for data
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

	for r := 1; r <= 25; r++ {
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
		quote.UpdRow(tbl)()
		quote.ClrBold(tbl)()
		quote.UpdOverviewRow(overviewTbl)()
		quote.ClrOverviewBold(overviewTbl)()
	}

	return tbl
}

// Stream connects to websocket server and streams price quotes
func (m *Market) Stream(dataCh chan cq.Quoter) error {
	if err := connectWS(m, dataCh); err != nil {
		return err
	}

	return nil
}

// formats pair with uppercase letters separated by "/"
func fmtID(s string) string {
	temp := strings.Split(s, "")
	b := strings.Builder{}

	for i := 0; i < 3; i++ {
		b.WriteString(temp[i])
	}
	b.WriteString("/")
	for i := 3; i < len(temp); i++ {
		b.WriteString(temp[i])
	}

	return b.String()
}
