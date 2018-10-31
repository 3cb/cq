package gemini

import (
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/overview"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
)

// Quote contains most recent data for each crypto currency pair
type Quote struct {
	Symbol string

	ID     string
	Price  string
	Size   string
	Bid    string
	Ask    string
	Volume string
}

// MarketID returns the name of market as a string
func (quote Quote) MarketID() string {
	return "gemini"
}

// PairID returns the pair name as all caps with "-" separator
func (quote Quote) PairID() string {
	return quote.ID
}

// findTblRow uses the pair ID to determine the quote's table row
// Returns an int
func (quote Quote) findTblRow() int {
	switch quote.ID {
	case "BTC-USD":
		return 2
	case "ETH-USD":
		return 5
	case "ETH-BTC":
		return 7
	case "LTC-USD":
		return 10
	case "LTC-BTC":
		return 12
	case "ZEC-USD":
		return 15
	// case "ZEC-BTC":
	default:
		return 17
	}
}

// UpdRow refreshes table with new data from websocket message
func (quote Quote) UpdRow(table *tview.Table) func() {
	return func() {
		row := quote.findTblRow()
		// delta, color := cq.FmtDelta(quote.Price, quote.Open)

		table.GetCell(row, 0).
			SetText(cq.FmtPair(quote.ID)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 1).
			SetText(cq.FmtPrice(quote.Price)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		// table.GetCell(row, 2).
		// 	SetText(delta).
		// 	SetTextColor(color).
		// 	SetAttributes(tcell.AttrBold)
		table.GetCell(row, 2).
			SetText(cq.FmtSize(quote.Size)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 3).
			SetText(cq.FmtPrice(quote.Bid)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 4).
			SetText(cq.FmtPrice(quote.Ask)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		// table.GetCell(row, 6).
		// 	SetText(cq.FmtPrice(quote.Low)).
		// 	SetTextColor(color).
		// 	SetAttributes(tcell.AttrBold)
		// table.GetCell(row, 7).
		// 	SetText(cq.FmtPrice(quote.High)).
		// 	SetTextColor(color).
		// 	SetAttributes(tcell.AttrBold)
		table.GetCell(row, 5).
			SetText(cq.FmtVolume(quote.Volume)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
	}
}

// ClrBold resets "Price" cell's attributes to remove bold font
func (quote Quote) ClrBold(table *tview.Table) func() {
	return func() {
		row := quote.findTblRow()

		for col := 0; col <= 5; col++ {
			table.GetCell(row, col).
				SetAttributes(tcell.AttrNone)
		}
	}
}

// UpdOverviewRow resets price quote in overview display
func (quote Quote) UpdOverviewRow(table *tview.Table) func() {
	return func() {
		row := overview.FindRow(quote)
		col := overview.FindColumn(quote)
		// _, color := cq.FmtDelta(quote.Price, quote.Open)

		table.GetCell(row, col).
			SetText(cq.FmtPrice(quote.Price)).
			// SetTextColor(color).
			SetAttributes(tcell.AttrBold)
	}
}

// ClrOverviewBold removes bold font from Price cells in overview display
func (quote Quote) ClrOverviewBold(table *tview.Table) func() {
	return func() {
		row := overview.FindRow(quote)
		col := overview.FindColumn(quote)

		table.GetCell(row, col).
			SetAttributes(tcell.AttrNone)
	}
}
