package bitfinex

import (
	"strconv"
	"strings"

	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type table struct {
	*tview.Table
}

// NewTable returns an instance of tview.Table formatted for bitfinex ready for data
func (m *Market) NewTable(cfg cq.Config) cq.Quoter {
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

	t := &table{
		tview.NewTable().
			SetBorders(true).
			SetBordersColor(cfg.Theme.BorderColor),
	}

	t.SetBackgroundColor(cfg.Theme.BackgroundColor)

	for i, header := range headers {
		t.SetCell(0, i, tview.NewTableCell(header).
			SetTextColor(cfg.Theme.HeaderTextColor).
			SetAlign(tview.AlignRight))
	}

	for r := 1; r < len(m.pairs)+1; r++ {
		for c := 0; c < len(headers); c++ {
			t.SetCell(r, c, tview.NewTableCell("-").
				SetTextColor(cfg.Theme.TextColor).
				SetAlign(tview.AlignRight))
		}
	}

	// handle errors here ***************
	m.getSnapshot()

	data := m.GetQuotes()

	for _, quote := range data {
		msg := cq.UpdateMsg{
			Quote:   quote,
			IsTrade: false,
			Flash:   false,
		}
		t.InsertQuote(msg, cfg)
	}

	return t
}

// InsertQuote updates table with new data
func (t *table) InsertQuote(msg cq.UpdateMsg, cfg cq.Config) {
	q := msg.Quote
	row := findTblRow(q)
	delta, color := FmtDelta(q.ChangePerc, cfg)

	t.GetCell(row, 0).
		SetText(q.ID).
		SetTextColor(color)
	t.GetCell(row, 1).
		SetText(cq.FmtPrice(q.Price)).
		SetTextColor(color)

	// if update is trade, set flash attribute on price cell
	if msg.IsTrade {
		if msg.Flash == true {
			t.GetCell(row, 1).
				SetAttributes(cfg.CellFlash)
		} else {
			t.GetCell(row, 1).
				SetAttributes(tcell.AttrNone)
		}
	}

	t.GetCell(row, 2).
		SetText(delta).
		SetTextColor(color)
	t.GetCell(row, 3).
		SetText(cq.FmtSize(q.Size)).
		SetTextColor(color)
	t.GetCell(row, 4).
		SetText(cq.FmtPrice(q.Bid)).
		SetTextColor(color)
	t.GetCell(row, 5).
		SetText(cq.FmtPrice(q.Ask)).
		SetTextColor(color)
	t.GetCell(row, 6).
		SetText(cq.FmtPrice(q.Low)).
		SetTextColor(color)
	t.GetCell(row, 7).
		SetText(cq.FmtPrice(q.High)).
		SetTextColor(color)
	t.GetCell(row, 8).
		SetText(cq.FmtVolume(q.Volume)).
		SetTextColor(color)
}

// FmtDelta is used for bitfinex exchange because it
// does not include an open price so cq.FmtDelta can't
// be used
func FmtDelta(change string, cfg cq.Config) (string, tcell.Color) {
	delta, err := strconv.ParseFloat(change, 64)
	if err != nil {
		return "no data", cfg.Theme.TextColor
	}
	delta *= 100
	d := strconv.FormatFloat(delta, 'f', 2, 64)
	spc := 9 - len(d) - 1 // extra "1" acounts for "%" character
	b := strings.Builder{}
	for i := 0; i < spc; i++ {
		b.WriteString(" ")
	}
	b.WriteString(d)
	b.WriteString("%")
	if delta > 0 {
		return b.String(), cfg.Theme.TextColorUp
	}
	return b.String(), cfg.Theme.TextColorDown
}

// FindTblRow uses the pair ID to determine the quote's table row
// Returns an int
func findTblRow(q cq.Quote) int {
	switch q.ID {
	case "BTC/USD":
		return 1
	case "BTC/EUR":
		return 2
	case "BTC/GBP":
		return 3
	case "BTC/JPY":
		return 4
	case "BCH/USD":
		return 5
	case "BCH/BTC":
		return 6
	case "ETH/USD":
		return 7
	case "ETH/BTC":
		return 8
	case "ETH/EUR":
		return 9
	case "ETH/GBP":
		return 10
	case "ETH/JPY":
		return 11
	case "LTC/USD":
		return 12
	case "LTC/BTC":
		return 13
	case "ZEC/USD":
		return 14
	case "ZEC/BTC":
		return 15
	case "ZRX/USD":
		return 16
	// case "ZRX/BTC":
	default:
		return 17
	}
}
