package cq

import (
	"github.com/3cb/tview"
)

// Exchange interface allows caller to get http snapshot price quotes,
// stream live data via websocket, create tables for display in gui
type Exchange interface {
	// GetPairs returns a slice of strings which is a list of products traded
	GetPairs() []string
	// GetSnapshot makes http requests to prime display table with data
	GetSnapshot() []error
	// Table returns display table with initial data
	Table() *tview.Table
	// PrimeOverview takes all price quotes from exchange to seed overview display
	PrimeTables(chan Quoter)
	// Stream launches goroutine to stream price data to display table
	Stream(chan Quoter) error
	// GetQuotes returns a map of price quotes for each product pair
	GetQuotes() map[string]Quoter
}
