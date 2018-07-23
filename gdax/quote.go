package gdax

import (
	"github.com/3cb/cq/cq"
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
	return "gdax"
}

// PairID returns name of product pair as a string
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
	case "BCH-USD":
		return 9
	case "BCH-BTC":
		return 11
	case "BCH-EUR":
		return 13
	case "ETH-USD":
		return 16
	case "ETH-BTC":
		return 18
	case "ETH-EUR":
		return 20
	case "LTC-USD":
		return 23
	case "LTC-BTC":
		return 25
	// case "LTC-EUR":
	default:
		return 27
	}
}

// SetRow sets tview.Table cells with data from an instance of Quote
func (quote Quote) SetRow(table *tview.Table) {
	r := quote.FindTblRow()
	delta, color := cq.FmtDelta(quote.Price, quote.Open)

	table.SetCell(r, 0, tview.NewTableCell(cq.FmtPair(quote.ID)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 1, tview.NewTableCell(cq.FmtPrice(quote.Price)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 2, tview.NewTableCell(delta).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 3, tview.NewTableCell(cq.FmtSize(quote.Size)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 4, tview.NewTableCell(cq.FmtPrice(quote.Bid)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 5, tview.NewTableCell(cq.FmtPrice(quote.Ask)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 6, tview.NewTableCell(cq.FmtPrice(quote.Low)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 7, tview.NewTableCell(cq.FmtPrice(quote.High)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
	table.SetCell(r, 8, tview.NewTableCell(cq.FmtVolume(quote.Volume)).
		SetTextColor(color).
		SetAlign(tview.AlignRight))
}

// UpdRow refreshes table with new data from websocket message
func (quote Quote) UpdRow(table *tview.Table) {
	row := quote.FindTblRow()
	delta, color := cq.FmtDelta(quote.Price, quote.Open)

	table.GetCell(row, 0).
		SetText(cq.FmtPair(quote.ID)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 1).
		SetText(cq.FmtPrice(quote.Price)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 2).
		SetText(delta).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 3).
		SetText(cq.FmtSize(quote.Size)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 4).
		SetText(cq.FmtPrice(quote.Bid)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 5).
		SetText(cq.FmtPrice(quote.Ask)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 6).
		SetText(cq.FmtPrice(quote.Low)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 7).
		SetText(cq.FmtPrice(quote.High)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
	table.GetCell(row, 8).
		SetText(cq.FmtVolume(quote.Volume)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
}

// ClrBold resets "Price" cell's attributes to remove bold font
func (quote Quote) ClrBold(table *tview.Table) {
	row := quote.FindTblRow()

	table.GetCell(row, 1).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 0).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 1).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 2).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 3).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 4).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 5).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 6).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 7).
		SetAttributes(tcell.AttrNone)
	table.GetCell(row, 8).
		SetAttributes(tcell.AttrNone)
}

// FindOverviewRow returns table row as integer
func (quote Quote) FindOverviewRow() int {
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
		return 10
	case "BCH-BTC":
		return 12
	case "BCH-EUR":
		return 14
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
		return 26
	case "LTC-BTC":
		return 28
	// case "LTC-EUR":
	default:
		return 30
	}
}

// UpdOverviewRow resets price quote in overview display
func (quote Quote) UpdOverviewRow(table *tview.Table) {
	row := quote.FindOverviewRow()
	_, color := cq.FmtDelta(quote.Price, quote.Open)

	table.GetCell(row, 1).
		SetText(cq.FmtPrice(quote.Price)).
		SetTextColor(color).
		SetAttributes(tcell.AttrBold)
}
