package hitbtc

import (
	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type table struct {
	*tview.Table
}

// NewTable returns an instance of tview.Table formatted for hitbtc ready for data
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

	for r := 1; r < len(m.symbols)+1; r++ {
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

func (t *table) InsertQuote(msg cq.UpdateMsg, cfg cq.Config) {

	q := msg.Quote
	row := findTblRow(q)
	delta, color := cq.FmtDelta(q.Price, q.Open, cfg)

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

// returns appropriate row as int
func findTblRow(q cq.Quote) int {
	switch q.ID {
	case "BTC/USD":
		return 1
	case "BCH/USD":
		return 2
	case "ETH/USD":
		return 3
	case "ETH/BTC":
		return 4
	case "LTC/USD":
		return 5
	case "LTC/BTC":
		return 6
	case "ZEC/USD":
		return 7
	case "ZEC/BTC":
		return 8
	case "ZRX/USD":
		return 9
	// case "ZRX/BTC":
	default:
		return 10
	}
}
