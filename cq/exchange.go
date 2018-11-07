package cq

import (
	"github.com/rivo/tview"
)

// Exchange interface allows caller to get http snapshot price quotes,
// stream live data via websocket, create tables for display in gui
type Exchange interface {
	// GetIDs returns a slice of trading pairs formatted with "/" (i.e., BTC/USD)
	GetIDs() []string

	// Table returns display table with initial data
	Table(*tview.Table) *tview.Table

	// Stream launches goroutine to stream price data to display table
	Stream(chan<- TimerMsg) error
}
