package overview

import (
	"github.com/3cb/cq/bitfinex"
	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// table wraps a tview.Table and provides extra behavior
type table struct {
	*tview.Table
}

// NewTable creates table comparing prices from different crypto exchanges
func NewTable(cfg cq.Config) cq.Quoter {
	rowLbl := []string{
		"BTC/USD",
		"BTC/EUR",
		"BTC/GBP",
		"BTC/JPY",
		"BCH/USD",
		"BCH/BTC",
		"BCH/EUR",
		"ETH/USD",
		"ETH/BTC",
		"ETH/EUR",
		"ETH/GBP",
		"ETH/JPY",
		"LTC/USD",
		"LTC/BTC",
		"LTC/EUR",
		"ZEC/USD",
		"ZEC/BTC",
		"ZRX/USD",
		"ZRX/BTC",
	}
	t := &table{
		tview.NewTable().
			SetBorders(true).
			SetBordersColor(cfg.Theme.BorderColor),
	}

	t.SetBackgroundColor(cfg.Theme.BackgroundColor)

	t.SetCell(0, 0, tview.NewTableCell("Pair").
		SetTextColor(cfg.Theme.HeaderTextColor).
		SetAlign(tview.AlignRight))

	t.SetCell(0, 1, tview.NewTableCell("  Coinbase").
		SetTextColor(cfg.Theme.HeaderTextColor).
		SetAlign(tview.AlignRight))

	t.SetCell(0, 2, tview.NewTableCell("  Bitfinex").
		SetTextColor(cfg.Theme.HeaderTextColor).
		SetAlign(tview.AlignRight))

	t.SetCell(0, 3, tview.NewTableCell("    HitBTC").
		SetTextColor(cfg.Theme.HeaderTextColor).
		SetAlign(tview.AlignRight))

	for i, pair := range rowLbl {
		row := i + 1
		t.SetCell(row, 0, tview.NewTableCell(pair).
			SetTextColor(cfg.Theme.TextColor).
			SetAlign(tview.AlignRight))

		t.SetCell(row, 1, tview.NewTableCell("-").
			SetTextColor(cfg.Theme.TextColor).
			SetAlign(tview.AlignRight))

		t.SetCell(row, 2, tview.NewTableCell("-").
			SetTextColor(cfg.Theme.TextColor).
			SetAlign(tview.AlignRight))

		t.SetCell(row, 3, tview.NewTableCell("-").
			SetTextColor(cfg.Theme.TextColor).
			SetAlign(tview.AlignRight))
	}

	return t
}

// InsertQuote updates display of market quotes with new data from UpdateMsg
func (t *table) InsertQuote(msg cq.UpdateMsg, cfg cq.Config) {
	row := FindRow(msg.Quote)
	col := FindColumn(msg.Quote)
	var color tcell.Color

	if msg.Quote.MarketID == "bitfinex" {
		_, color = bitfinex.FmtDelta(msg.Quote.ChangePerc, cfg)
	} else {
		_, color = cq.FmtDelta(msg.Quote.Price, msg.Quote.Open, cfg)
	}

	var mask tcell.AttrMask
	if msg.Flash == true {
		mask = cfg.CellFlash
	} else {
		mask = tcell.AttrNone
	}

	t.GetCell(row, col).
		SetText(cq.FmtPrice(msg.Quote.Price)).
		SetTextColor(color).
		SetAttributes(mask)
}

// FindRow uses pair string to find row in table
func FindRow(q cq.Quote) int {
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
	case "BCH/EUR":
		return 7
	case "ETH/USD":
		return 8
	case "ETH/BTC":
		return 9
	case "ETH/EUR":
		return 10
	case "ETH/GBP":
		return 11
	case "ETH/JPY":
		return 12
	case "LTC/USD":
		return 13
	case "LTC/BTC":
		return 14
	case "LTC/EUR":
		return 15
	case "ZEC/USD":
		return 16
	case "ZEC/BTC":
		return 17
	case "ZRX/USD":
		return 18
	// case "ZRX/BTC":
	default:
		return 19
	}
}

// FindColumn uses MarketID string to find column in table
func FindColumn(q cq.Quote) int {
	switch q.MarketID {
	case "coinbase":
		return 1
	case "bitfinex":
		return 2
	// case "hitbtc":
	default:
		return 3
	}
}
