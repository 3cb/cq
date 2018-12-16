package cq

import tview "github.com/rivo/tview"

// Quoter defines methods for all display tables in cq application necessary
// to display price quotes for supported exchanges
// InsertQuote returns a closure
// which update the tview.Table with new data.
// The closure is then passed on to the
// tview.Application via app.QueueUpdateDraw()
type Quoter interface {
	// Primitive interface is embedded to ease use of Quoter in main.go
	tview.Primitive

	// InsertQuote updates table values as well as styling (i.e., tcell.AttrMask)
	InsertQuote(UpdateMsg, Config)
}
