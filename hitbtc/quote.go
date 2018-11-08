package hitbtc

import (
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/overview"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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

func (q Quote) InsertTicker(table *tview.Table) func() {
	return func() {
		row := q.findTblRow()
		delta, color := cq.FmtDelta(q.Price, q.Open)

		table.GetCell(row, 0).
			SetText(q.ID).
			SetTextColor(color)
		table.GetCell(row, 1).
			SetText(cq.FmtPrice(q.Price)).
			SetTextColor(color)
		table.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color)
		table.GetCell(row, 3).
			SetText(cq.FmtSize(q.Size)).
			SetTextColor(color)
		table.GetCell(row, 4).
			SetText(cq.FmtPrice(q.Bid)).
			SetTextColor(color)
		table.GetCell(row, 5).
			SetText(cq.FmtPrice(q.Ask)).
			SetTextColor(color)
		table.GetCell(row, 6).
			SetText(cq.FmtPrice(q.Low)).
			SetTextColor(color)
		table.GetCell(row, 7).
			SetText(cq.FmtPrice(q.High)).
			SetTextColor(color)
		table.GetCell(row, 8).
			SetText(cq.FmtVolume(q.Volume)).
			SetTextColor(color)
	}
}

func (q Quote) InsertTrade(oTable *tview.Table, table *tview.Table, attr tcell.AttrMask) func() {
	return func() {
		// update exchange table
		row := q.findTblRow()
		delta, color := cq.FmtDelta(q.Price, q.Open)

		table.GetCell(row, 0).
			SetText(q.ID).
			SetTextColor(color)
		table.GetCell(row, 1).
			SetText(cq.FmtPrice(q.Price)).
			SetTextColor(color).
			SetAttributes(attr)
		table.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color)
		table.GetCell(row, 3).
			SetText(cq.FmtSize(q.Size)).
			SetTextColor(color)
		table.GetCell(row, 4).
			SetText(cq.FmtPrice(q.Bid)).
			SetTextColor(color)
		table.GetCell(row, 5).
			SetText(cq.FmtPrice(q.Ask)).
			SetTextColor(color)
		table.GetCell(row, 6).
			SetText(cq.FmtPrice(q.Low)).
			SetTextColor(color)
		table.GetCell(row, 7).
			SetText(cq.FmtPrice(q.High)).
			SetTextColor(color)
		table.GetCell(row, 8).
			SetText(cq.FmtVolume(q.Volume)).
			SetTextColor(color)

		// update overview table
		row = overview.FindRow(q)
		col := overview.FindColumn(q)
		_, color = cq.FmtDelta(q.Price, q.Open)

		oTable.GetCell(row, col).
			SetText(cq.FmtPrice(q.Price)).
			SetTextColor(color).
			SetAttributes(attr)
	}
}
