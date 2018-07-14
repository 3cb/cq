package cq

import (
	"github.com/rivo/tview"
)

// Exchange interface allows caller to get http snapshot price quotes,
// stream live data via websocket, create tables for display in gui
type Exchange interface {
	Snapshot() []error
	Stream(chan Setter) error
	Table() *tview.Table
	// Print()
}

// MarketComp contains most current price for each pair from each exchange
type MarketComp struct {
	GDAX     string
	Gemini   string
	Bitfinex string
}
