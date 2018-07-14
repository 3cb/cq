package cq

import (
	"github.com/rivo/tview"
)

// Setter is an interface that market quotes implement to initialize
// and update data in gui table
type Setter interface {
	MarketID() string

	// SetRow(*tview.Table)
	UpdRow(*tview.Table)
	ClrBold(*tview.Table)

	PrimeOverview(chan Setter)
	UpdOverviewRow(*tview.Table)
	// ClrOverviewBold(*tview.Table)
}
