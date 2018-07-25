package cq

import (
	"github.com/rivo/tview"
)

// Quoter is an interface that market quotes implement to initialize
// and update data in gui table
type Quoter interface {
	MarketID() string
	PairID() string

	UpdRow(*tview.Table)
	ClrBold(*tview.Table)

	UpdOverviewRow(*tview.Table)
	ClrOverviewBold(*tview.Table)
}
