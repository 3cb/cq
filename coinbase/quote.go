package coinbase

import (
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/overview"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Quote contains most recent data for each crypto currency pair
// Trailing comments denote which http request or websocket stream the data comes from
// getTrades: https://docs.gdax.com/#get-trades
// match: https://docs.gdax.com/#the-code-classprettyprintfullcode-channel
// ticker: https://docs.gdax.com/#the-code-classprettyprinttickercode-channel
// *** GDAX API documentation for websocket ticker channel does not show all available fields as of 2/11/2018
type Quote struct {
	Type string `json:"type"` // used to filter websocket messages

	ID    string `json:"product_id"`
	Price string `json:"price"` // getTrades/match
	Size  string `json:"size"`  // getTrades/match

	Bid string `json:"best_bid"` // getTicker/ticker
	Ask string `json:"best_ask"` // getTicker/ticker

	High   string `json:"high_24h"`   // getStats/ticker
	Low    string `json:"low_24h"`    // getStats/ticker
	Open   string `json:"open_24h"`   // getStats/ticker
	Volume string `json:"volume_24h"` // getStats/ticker
}

// MarketID returns the name of market as a string
func (quote Quote) MarketID() string {
	return "coinbase"
}

// PairID returns name of product pair as a string
func (quote Quote) PairID() string {
	return quote.ID
}

// findTblRow uses the pair ID to determine the quote's table row
// Returns an int
func (quote Quote) findTblRow() int {
	switch quote.ID {
	case "BTC/USD":
		return 1
	case "BTC/EUR":
		return 2
	case "BTC/GBP":
		return 3
	case "BCH/USD":
		return 4
	case "BCH/BTC":
		return 5
	case "BCH/EUR":
		return 6
	case "ETH/USD":
		return 7
	case "ETH/BTC":
		return 8
	case "ETH/EUR":
		return 9
	case "ETH/GBP":
		return 10
	case "LTC/USD":
		return 11
	case "LTC/BTC":
		return 12
	case "LTC/EUR":
		return 13
	case "ZRX/USD":
		return 14
	// case "ZRX/BTC":
	default:
		return 15
	}
}

func (quote Quote) InsertTicker(table *tview.Table) func() {
	return func() {
		row := quote.findTblRow()
		delta, color := cq.FmtDelta(quote.Price, quote.Open)

		table.GetCell(row, 0).
			SetText(quote.ID).
			SetTextColor(color)
		table.GetCell(row, 1).
			SetText(cq.FmtPrice(quote.Price)).
			SetTextColor(color)
		table.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color)
		table.GetCell(row, 3).
			SetText(cq.FmtSize(quote.Size)).
			SetTextColor(color)
		table.GetCell(row, 4).
			SetText(cq.FmtPrice(quote.Bid)).
			SetTextColor(color)
		table.GetCell(row, 5).
			SetText(cq.FmtPrice(quote.Ask)).
			SetTextColor(color)
		table.GetCell(row, 6).
			SetText(cq.FmtPrice(quote.Low)).
			SetTextColor(color)
		table.GetCell(row, 7).
			SetText(cq.FmtPrice(quote.High)).
			SetTextColor(color)
		table.GetCell(row, 8).
			SetText(cq.FmtVolume(quote.Volume)).
			SetTextColor(color)
	}
}

func (quote Quote) InsertTrade(oTable *tview.Table, table *tview.Table, attr tcell.AttrMask) func() {
	return func() {
		row := quote.findTblRow()
		delta, color := cq.FmtDelta(quote.Price, quote.Open)

		table.GetCell(row, 0).
			SetText(quote.ID).
			SetTextColor(color)
		table.GetCell(row, 1).
			SetText(cq.FmtPrice(quote.Price)).
			SetTextColor(color).
			SetAttributes(attr)
		table.GetCell(row, 2).
			SetText(delta).
			SetTextColor(color)
		table.GetCell(row, 3).
			SetText(cq.FmtSize(quote.Size)).
			SetTextColor(color)
		table.GetCell(row, 4).
			SetText(cq.FmtPrice(quote.Bid)).
			SetTextColor(color)
		table.GetCell(row, 5).
			SetText(cq.FmtPrice(quote.Ask)).
			SetTextColor(color)
		table.GetCell(row, 6).
			SetText(cq.FmtPrice(quote.Low)).
			SetTextColor(color)
		table.GetCell(row, 7).
			SetText(cq.FmtPrice(quote.High)).
			SetTextColor(color)
		table.GetCell(row, 8).
			SetText(cq.FmtVolume(quote.Volume)).
			SetTextColor(color)

		// update overview table
		row = overview.FindRow(quote)
		col := overview.FindColumn(quote)
		_, color = cq.FmtDelta(quote.Price, quote.Open)

		oTable.GetCell(row, col).
			SetText(cq.FmtPrice(quote.Price)).
			SetTextColor(color).
			SetAttributes(attr)
	}
}
