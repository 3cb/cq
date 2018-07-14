package overview

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Table creates table comparing prices from different crypto exchanges
func Table(exchanges map[string]cq.Exchange) *tview.Table {
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
	}

	table := tview.NewTable().
		SetBorders(false)

	table.SetCell(0, 1, tview.NewTableCell("              GDAX").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 2, tview.NewTableCell("          Bitfinex").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 3, tview.NewTableCell("            Gemini").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	for i, pair := range rowLbl {
		table.SetCell(2*(i+1), 0, tview.NewTableCell(pair).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))
	}

	return table
}
