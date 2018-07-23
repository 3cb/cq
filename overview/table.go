package overview

import (
	"github.com/3cb/cq/cq"
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
		row := 2 * (i + 1)
		table.SetCell(row, 0, tview.NewTableCell(pair).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))

		table.SetCell(row, 1, tview.NewTableCell("").
			SetAlign(tview.AlignRight))

		table.SetCell(row, 2, tview.NewTableCell("").
			SetAlign(tview.AlignRight))

		table.SetCell(row, 3, tview.NewTableCell("").
			SetAlign(tview.AlignRight))

	}

	return table
}
