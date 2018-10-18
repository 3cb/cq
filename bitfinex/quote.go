package bitfinex

import (
	"strconv"
	"strings"

	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/overview"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
)

// Quote represents data for a crypto trading pair on bitfinex exchange
// data comes from websocket messages in array format
type Quote struct {
	Symbol string

	ID         string
	Price      float64
	Change     float64
	ChangePerc float64
	Size       float64

	Bid float64
	Ask float64

	Low    float64
	High   float64
	Open   float64
	Volume float64
}

// MarketID returns the name of market as a string
func (quote Quote) MarketID() string {
	return "bitfinex"
}

// PairID returns the name of product pair as a string
func (quote Quote) PairID() string {
	return quote.ID
}

// FindTblRow uses the pair ID to determine the quote's table row
// Returns an int
func (quote Quote) FindTblRow() int {
	switch quote.ID {
	case "BTC-USD":
		return 2
	case "BTC-EUR":
		return 4
	case "BTC-GBP":
		return 6
	case "BTC-JPY":
		return 8
	case "BCH-USD":
		return 11
	case "BCH-BTC":
		return 13
	case "ETH-USD":
		return 16
	case "ETH-BTC":
		return 18
	case "ETH-EUR":
		return 20
	case "ETH-GBP":
		return 22
	case "ETH-JPY":
		return 24
	case "LTC-USD":
		return 27
	// case "LTC-BTC":
	default:
		return 29
	}
}

// UpdRow refreshes table with new data from websocket message
func (quote Quote) UpdRow(table *tview.Table) {
	var color tcell.Color

	if quote.ChangePerc >= 0 {
		color = tcell.ColorGreen
	} else {
		color = tcell.ColorRed
	}
	price := strconv.FormatFloat(quote.Price, 'f', -1, 64)

	delta := formatDelta(quote.ChangePerc)

	size := strconv.FormatFloat(quote.Size, 'f', -1, 64)
	bid := strconv.FormatFloat(quote.Bid, 'f', -1, 64)
	ask := strconv.FormatFloat(quote.Ask, 'f', -1, 64)
	low := strconv.FormatFloat(quote.Low, 'f', -1, 64)
	high := strconv.FormatFloat(quote.High, 'f', -1, 64)
	vol := strconv.FormatFloat(quote.Volume, 'f', -1, 64)

	row := quote.FindTblRow()

	table.GetCell(row, 0).
		SetText(cq.FmtPair(quote.ID)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 1).
		SetText(cq.FmtPrice(price)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 2).
		SetText(delta).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 3).
		SetText(cq.FmtSize(size)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 4).
		SetText(cq.FmtPrice(bid)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 5).
		SetText(cq.FmtPrice(ask)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 6).
		SetText(cq.FmtPrice(low)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 7).
		SetText(cq.FmtPrice(high)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 8).
		SetText(cq.FmtVolume(vol)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
}

// ClrBold resets "Price" cell's attributes to remove bold font
func (quote Quote) ClrBold(table *tview.Table) {
	row := quote.FindTblRow()

	for col := 0; col <= 8; col++ {
		table.GetCell(row, col).
			SetAttributes(tcell.AttrNone)
	}
}

// UpdOverviewRow resets price quote in overview display
func (quote Quote) UpdOverviewRow(table *tview.Table) {
	var color tcell.Color

	row := overview.FindRow(quote)
	if quote.Change >= 0 {
		color = tcell.ColorGreen
	} else {
		color = tcell.ColorRed
	}

	price := strconv.FormatFloat(quote.Price, 'f', -1, 64)

	table.GetCell(row, 2).
		SetText(cq.FmtPrice(price)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
}

// ClrOverviewBold removes bold font from Price cells in overview display
func (quote Quote) ClrOverviewBold(table *tview.Table) {
	row := overview.FindRow(quote)

	table.GetCell(row, 2).
		SetAttributes(tcell.AttrNone)
}

func formatDelta(change float64) string {
	change *= 100
	delta := strconv.FormatFloat(change, 'f', 2, 64)
	spc := 9 - len(delta) - 1 // extra "1" acounts for "%" character
	b := strings.Builder{}
	for i := 0; i < spc; i++ {
		b.WriteString(" ")
	}
	b.WriteString(delta)
	b.WriteString("%")
	return b.String()
}
