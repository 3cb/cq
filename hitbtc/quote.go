package hitbtc

import (
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/overview"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
)

// Quote holds data for each trading pair
type Quote struct {
	Symbol string
	ID     string // formatted with "/"

	Price string
	// Change     string
	// ChangePerc string
	Size   string
	Bid    string
	Ask    string
	Open   string
	Low    string
	High   string
	Volume string
}

// MarketID returns the name of market as a string
func (q Quote) MarketID() string {
	return "hitbtc"
}

// PairID returns the trading pair as a string formatted with "/"
func (q Quote) PairID() string {
	return q.ID
}

// returns appropriate row as int
func (q Quote) findTblRow() int {
	switch q.ID {
	case "BTC/USD":
		return 2
	case "BCH/USD":
		return 5
	case "ETH/USD":
		return 8
	case "ETH/BTC":
		return 10
	case "LTC/USD":
		return 13
	case "LTC/BTC":
		return 15
	case "ZEC/USD":
		return 18
	case "ZEC/BTC":
		return 20
	case "ZRX/USD":
		return 23
	// case "ZRX/BTC":
	default:
		return 25
	}
}

// UpdRow refreshes table with new data from websocket message
func (q Quote) UpdRow(table *tview.Table) func() {
	return func() {
		row := q.findTblRow()
		delta, color := cq.FmtDelta(q.Price, q.Open)

		table.GetCell(row, 0).
			SetText(q.ID).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 1).
			SetText(cq.FmtPrice(q.Price)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 3).
			SetText(cq.FmtSize(q.Size)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 4).
			SetText(cq.FmtPrice(q.Bid)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 5).
			SetText(cq.FmtPrice(q.Ask)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 6).
			SetText(cq.FmtPrice(q.Low)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 7).
			SetText(cq.FmtPrice(q.High)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
		table.GetCell(row, 8).
			SetText(cq.FmtVolume(q.Volume)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
	}
}

// ClrBold resets "Price" cell's attributes to remove bold font
func (q Quote) ClrBold(table *tview.Table) func() {
	return func() {
		row := q.findTblRow()

		for col := 0; col <= 8; col++ {
			table.GetCell(row, col).
				SetAttributes(tcell.AttrNone)
		}
	}
}

// UpdOverviewRow resets price quote in overview display
func (q Quote) UpdOverviewRow(table *tview.Table) func() {
	return func() {
		row := overview.FindRow(q)
		col := overview.FindColumn(q)
		_, color := cq.FmtDelta(q.Price, q.Open)

		table.GetCell(row, col).
			SetText(cq.FmtPrice(q.Price)).
			SetTextColor(color).
			SetAttributes(tcell.AttrBold)
	}
}

// ClrOverviewBold removes bold font from Price cells in overview display
func (q Quote) ClrOverviewBold(table *tview.Table) func() {
	return func() {
		row := overview.FindRow(q)
		col := overview.FindColumn(q)

		table.GetCell(row, col).
			SetAttributes(tcell.AttrNone)
	}
}
