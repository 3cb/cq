package cq

import (
	"github.com/rivo/tview"
)

// Quoter is an interface that market quotes implement to initialize
// and update data in gui table
type Quoter interface {
	MarketID() string
	PairID() string

	// UpdRow updates exchange table row with new quote data and changes
	// font to bold to signal change in data
	UpdRow(*tview.Table) func()
	// ClrBold changes font of table row back to normal
	ClrBold(*tview.Table) func()

	UpdOverviewRow(*tview.Table) func()
	ClrOverviewBold(*tview.Table) func()
}
