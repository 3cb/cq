package cq

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// UpdateMsg carries quotes from TimerGroup event loop to main cq event loop
// UpdType and Flash fields allow event loop to set table fonts for quotes
type UpdateMsg struct {
	UpdType string // "trade" or "ticker"
	Flash   bool
	Quote   Quoter
}

// Quoter is an interface that market quotes implement to initialize
// and update data in gui table
type Quoter interface {
	// returns the name of the exchange as all lowercase string
	MarketID() string

	// returns name of trading pair all caps separated by "/" (e.g., BTC/USD)
	PairID() string

	// Insert new quote data into market and overview tables with formatting for flash
	InsertTrade(*tview.Table, *tview.Table, tcell.AttrMask) func()

	// Insert new quote data into market table without altering flash state
	InsertTicker(*tview.Table) func()
}
