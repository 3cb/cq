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
	MarketID() string
	PairID() string

	// UpdRow updates exchange table row with new quote data and changes
	// font to bold to signal change in data
	//
	// boolean parameter should be true for boldface quote
	// UpdRow(*tview.Table, tcell.AttrMask) func()

	// UpdOverviewRow(*tview.Table, tcell.AttrMask) func()

	TradeUpdate(*tview.Table, *tview.Table, tcell.AttrMask) func()
	TickerUpdate(*tview.Table) func()
}
