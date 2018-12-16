package coinbase

import (
	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type table struct {
	*tview.Table
}

// NewTable method uses market data to create and return an
// instance of tview.Table that implements cq.Table
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
	quote := msg.Quote
	row := findTblRow(quote)
	delta, color := cq.FmtDelta(quote.Price, quote.Open, cfg)

	t.GetCell(row, 0).
		SetText(quote.ID).
		SetTextColor(color)
	t.GetCell(row, 1).
		SetText(cq.FmtPrice(quote.Price)).
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
		SetText(cq.FmtSize(quote.Size)).
		SetTextColor(color)
	t.GetCell(row, 4).
		SetText(cq.FmtPrice(quote.Bid)).
		SetTextColor(color)
	t.GetCell(row, 5).
		SetText(cq.FmtPrice(quote.Ask)).
		SetTextColor(color)
	t.GetCell(row, 6).
		SetText(cq.FmtPrice(quote.Low)).
		SetTextColor(color)
	t.GetCell(row, 7).
		SetText(cq.FmtPrice(quote.High)).
		SetTextColor(color)
	t.GetCell(row, 8).
		SetText(cq.FmtVolume(quote.Volume)).
		SetTextColor(color)
}

// findTblRow uses the pair ID to determine the quote's table row
// Returns an int
func findTblRow(q cq.Quote) int {
	switch q.ID {
	case "BTC/USD":
		return 1
	case "BTC/EUR":
		return 2
	case "BTC/GBP":
		return 3
	case "BCH/USD":
		return 4
	case "BCH/BTC":
		return 5
	case "BCH/EUR":
		return 6
	case "ETH/USD":
		return 7
	case "ETH/BTC":
		return 8
	case "ETH/EUR":
		return 9
	case "ETH/GBP":
		return 10
	case "LTC/USD":
		return 11
	case "LTC/BTC":
		return 12
	case "LTC/EUR":
		return 13
	case "ZRX/USD":
		return 14
	// case "ZRX/BTC":
	default:
		return 15
	}
}
