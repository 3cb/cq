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
		SetBorders(true).
		SetBordersColor(tcell.ColorLightSlateGray)

	table.SetCell(0, 1, tview.NewTableCell("  Coinbase").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 2, tview.NewTableCell("  Bitfinex").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	table.SetCell(0, 3, tview.NewTableCell("    HitBTC").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignRight))

	for i, pair := range rowLbl {
		row := i + 1
		table.SetCell(row, 0, tview.NewTableCell(pair).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))

		table.SetCell(row, 1, tview.NewTableCell("-").
			SetAlign(tview.AlignRight))

		table.SetCell(row, 2, tview.NewTableCell("-").
			SetAlign(tview.AlignRight))

		table.SetCell(row, 3, tview.NewTableCell("-").
			SetAlign(tview.AlignRight))
	}

	return table
}

// FindRow uses pair string to find row in table
func FindRow(quote cq.Quoter) int {
	switch quote.PairID() {
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
func FindColumn(quote cq.Quoter) int {
	switch quote.MarketID() {
	case "coinbase":
		return 1
	case "bitfinex":
		return 2
	// case "hitbtc":
	default:
		return 3
	}
}
