package overview

import (
	"github.com/3cb/cq/cq"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Table creates table comparing prices from different crypto exchanges
func Table() *tview.Table {
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

	table := tview.NewTable().
		SetBorders(false)

	table.SetCell(0, 1, tview.NewTableCell("              GDAX").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 2, tview.NewTableCell("          Bitfinex").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 3, tview.NewTableCell("            HitBTC").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 4, tview.NewTableCell("            Gemini").
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

		table.SetCell(row, 4, tview.NewTableCell("").
			SetAlign(tview.AlignRight))
	}

	return table
}

// FindRow uses pair string to find row in table
func FindRow(quote cq.Quoter) int {
	switch quote.PairID() {
	case "BTC/USD":
		return 2
	case "BTC/EUR":
		return 4
	case "BTC/GBP":
		return 6
	case "BTC/JPY":
		return 8
	case "BCH/USD":
		return 10
	case "BCH/BTC":
		return 12
	case "BCH/EUR":
		return 14
	case "ETH/USD":
		return 16
	case "ETH/BTC":
		return 18
	case "ETH/EUR":
		return 20
	case "ETH/GBP":
		return 22
	case "ETH/JPY":
		return 24
	case "LTC/USD":
		return 26
	case "LTC/BTC":
		return 28
	case "LTC/EUR":
		return 30
	case "ZEC/USD":
		return 32
	case "ZEC/BTC":
		return 34
	case "ZRX/USD":
		return 36
	// case "ZRX/BTC":
	default:
		return 38
	}
}

// FindColumn uses MarketID string to find column in table
func FindColumn(quote cq.Quoter) int {
	switch quote.MarketID() {
	case "gdax":
		return 1
	case "bitfinex":
		return 2
	case "hitbtc":
		return 3
	// case "gemini":
	default:
		return 4
	}
}
