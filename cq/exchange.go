package cq

import "github.com/3cb/muttview"

// Exchange interface allows caller to get http snapshot price quotes,
// stream live data via websocket, create tables for display in gui
type Exchange interface {
	// Table returns display table with initial data
	Table(*tview.Table) *tview.Table

	// Stream launches goroutine to stream price data to display table
	Stream(chan Quoter) error
}
