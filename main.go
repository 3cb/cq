package main

import (
	"time"

	"github.com/gdamore/tcell"

	"github.com/3cb/cq/display"
	"github.com/3cb/cq/gdax"
	"github.com/rivo/tview"
)

func main() {
	exchanges := make(map[string]Exchange)
	exchanges["gdax"] = gdax.Init()
	// exchanges["gemini"] = gemini.Init()
	// exchanges["bitfinex"] = bitfinex.Init()

	// handle error slice here
	exchanges["gdax"].Snapshot()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack

	app := tview.NewApplication()

	menu := tview.NewList().
		// AddItem("Overview", "", '1', nil).
		AddItem("GDAX", "", '2', nil).
		AddItem("Gemini", "", '3', nil).
		AddItem("Bitfinex", "", '4', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	table := exchanges["gdax"].Table()

	body := tview.NewFlex().
		SetFullScreen(true).
		AddItem(menu, 20, 1, true).
		AddItem(table, 0, 1, false)

	gdaxStream := make(chan display.Setter, 100)

	go func() {
		exchanges["gdax"].Stream(gdaxStream)

		for {
			select {
			case upd := <-gdaxStream:
				upd.UpdRow(table)
				app.Draw()
				time.Sleep(100 * time.Millisecond)
				upd.ClrBold(table)
				app.Draw()
			}
		}
	}()

	if err := app.SetRoot(body, true).Run(); err != nil {
		panic(err)
	}
}

// Exchange interface allows caller to get http snapshot price quotes,
// stream live data via websocket, create tables for display in gui
type Exchange interface {
	Snapshot() []error
	Stream(chan display.Setter) error
	Table() *tview.Table
	// Print()
}
