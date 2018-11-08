package bitfinex

import (
	"strconv"
	"strings"

	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/overview"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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
func (quote Quote) findTblRow() int {
	switch quote.ID {
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
	case "ETH/USD":
		return 7
	case "ETH/BTC":
		return 8
	case "ETH/EUR":
		return 9
	case "ETH/GBP":
		return 10
	case "ETH/JPY":
		return 11
	case "LTC/USD":
		return 12
	case "LTC/BTC":
		return 13
	case "ZEC/USD":
		return 14
	case "ZEC/BTC":
		return 15
	case "ZRX/USD":
		return 16
	// case "ZRX/BTC":
	default:
		return 17
	}
}

func (quote Quote) TickerUpdate(tbl *tview.Table) func() {
	return func() {
		var color tcell.Color

		if quote.ChangePerc >= 0 {
			color = tcell.ColorGreen
		} else {
			color = tcell.ColorRed
		}
		price := strconv.FormatFloat(quote.Price, 'f', -1, 64)

		delta := fmtDelta(quote.ChangePerc)

		size := strconv.FormatFloat(quote.Size, 'f', -1, 64)
		bid := strconv.FormatFloat(quote.Bid, 'f', -1, 64)
		ask := strconv.FormatFloat(quote.Ask, 'f', -1, 64)
		low := strconv.FormatFloat(quote.Low, 'f', -1, 64)
		high := strconv.FormatFloat(quote.High, 'f', -1, 64)
		vol := strconv.FormatFloat(quote.Volume, 'f', -1, 64)

		row := quote.findTblRow()

		tbl.GetCell(row, 0).
			SetText(quote.ID).
			SetTextColor(color)
		tbl.GetCell(row, 1).
			SetText(cq.FmtPrice(price)).
			SetTextColor(color)
		tbl.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color)
		tbl.GetCell(row, 3).
			SetText(cq.FmtSize(size)).
			SetTextColor(color)
		tbl.GetCell(row, 4).
			SetText(cq.FmtPrice(bid)).
			SetTextColor(color)
		tbl.GetCell(row, 5).
			SetText(cq.FmtPrice(ask)).
			SetTextColor(color)
		tbl.GetCell(row, 6).
			SetText(cq.FmtPrice(low)).
			SetTextColor(color)
		tbl.GetCell(row, 7).
			SetText(cq.FmtPrice(high)).
			SetTextColor(color)
		tbl.GetCell(row, 8).
			SetText(cq.FmtVolume(vol)).
			SetTextColor(color)
	}
}

func (quote Quote) TradeUpdate(overviewTbl *tview.Table, tbl *tview.Table, attr tcell.AttrMask) func() {
	return func() {
		// update exchange table
		var color tcell.Color

		if quote.ChangePerc >= 0 {
			color = tcell.ColorGreen
		} else {
			color = tcell.ColorRed
		}
		price := strconv.FormatFloat(quote.Price, 'f', -1, 64)

		delta := fmtDelta(quote.ChangePerc)

		size := strconv.FormatFloat(quote.Size, 'f', -1, 64)
		bid := strconv.FormatFloat(quote.Bid, 'f', -1, 64)
		ask := strconv.FormatFloat(quote.Ask, 'f', -1, 64)
		low := strconv.FormatFloat(quote.Low, 'f', -1, 64)
		high := strconv.FormatFloat(quote.High, 'f', -1, 64)
		vol := strconv.FormatFloat(quote.Volume, 'f', -1, 64)

		row := quote.findTblRow()

		tbl.GetCell(row, 0).
			SetText(quote.ID).
			SetTextColor(color)
		tbl.GetCell(row, 1).
			SetText(cq.FmtPrice(price)).
			SetTextColor(color).
			SetAttributes(attr)
		tbl.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color)
		tbl.GetCell(row, 3).
			SetText(cq.FmtSize(size)).
			SetTextColor(color)
		tbl.GetCell(row, 4).
			SetText(cq.FmtPrice(bid)).
			SetTextColor(color)
		tbl.GetCell(row, 5).
			SetText(cq.FmtPrice(ask)).
			SetTextColor(color)
		tbl.GetCell(row, 6).
			SetText(cq.FmtPrice(low)).
			SetTextColor(color)
		tbl.GetCell(row, 7).
			SetText(cq.FmtPrice(high)).
			SetTextColor(color)
		tbl.GetCell(row, 8).
			SetText(cq.FmtVolume(vol)).
			SetTextColor(color)

		// update overview table
		row = overview.FindRow(quote)
		col := overview.FindColumn(quote)
		if quote.Change >= 0 {
			color = tcell.ColorGreen
		} else {
			color = tcell.ColorRed
		}

		price = strconv.FormatFloat(quote.Price, 'f', -1, 64)

		overviewTbl.GetCell(row, col).
			SetText(cq.FmtPrice(price)).
			SetTextColor(color).
			SetAttributes(attr)
	}
}

func fmtDelta(change float64) string {
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
