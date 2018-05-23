package gdax

import (
	"github.com/3cb/cq/display"
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

// FindTblRow uses the pair ID to determine the quote's table row
// Returns an int
func (quote Quote) FindTblRow() int {
	switch quote.ID {
	case "BTC-USD":
		return 1
	case "BTC-EUR":
		return 2
	case "BTC-GBP":
		return 3
	case "BCH-USD":
		return 4
	case "BCH-BTC":
		return 5
	case "BCH-EUR":
		return 6
	case "ETH-USD":
		return 7
	case "ETH-BTC":
		return 8
	case "ETH-EUR":
		return 9
	case "LTC-USD":
		return 10
	case "LTC-BTC":
		return 11
	// case "LTC-EUR":
	default:
		return 12
	}
}

// SetRow sets tview.Table cells with data from an instance of Quote
func (quote Quote) SetRow(table *tview.Table) {
	r := quote.FindTblRow()
	delta, color := display.FmtDelta(quote.Price, quote.Open)

	table.SetCell(r, 0, tview.NewTableCell(display.FmtPair(quote.ID)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 1, tview.NewTableCell(display.FmtPrice(quote.Price)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 2, tview.NewTableCell(delta).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 3, tview.NewTableCell(display.FmtSize(quote.Size)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 4, tview.NewTableCell(display.FmtPrice(quote.Bid)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 5, tview.NewTableCell(display.FmtPrice(quote.Ask)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 6, tview.NewTableCell(display.FmtPrice(quote.Low)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 7, tview.NewTableCell(display.FmtPrice(quote.High)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 8, tview.NewTableCell(display.FmtVolume(quote.Volume)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
}

// UpdRow refreshes table with new data from websocket message
func (quote Quote) UpdRow(table *tview.Table) {
	row := quote.FindTblRow()
	delta, color := display.FmtDelta(quote.Price, quote.Open)

	table.GetCell(row, 0).
		SetText(display.FmtPair(quote.ID)).
		SetTextColor(color)
	table.GetCell(row, 1).
		SetText(display.FmtPrice(quote.Price)).
		SetTextColor(color)
	table.GetCell(row, 2).
		SetText(delta).
		SetTextColor(color)
	table.GetCell(row, 3).
		SetText(display.FmtSize(quote.Size)).
		SetTextColor(color)
	table.GetCell(row, 4).
		SetText(display.FmtPrice(quote.Bid)).
		SetTextColor(color)
	table.GetCell(row, 5).
		SetText(display.FmtPrice(quote.Ask)).
		SetTextColor(color)
	table.GetCell(row, 6).
		SetText(display.FmtPrice(quote.Low)).
		SetTextColor(color)
	table.GetCell(row, 7).
		SetText(display.FmtPrice(quote.High)).
		SetTextColor(color)
	table.GetCell(row, 8).
		SetText(display.FmtVolume(quote.Volume)).
		SetTextColor(color)
}
