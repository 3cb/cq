package main

import (
	"time"

	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/gdax"
	"github.com/3cb/cq/overview"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	exchanges := make(map[string]cq.Exchange)
	exchanges["gdax"] = gdax.Init()
	// exchanges["gemini"] = gemini.Init()
	// exchanges["bitfinex"] = bitfinex.Init()

	showOverview := true

	// handle error slice here
	exchanges["gdax"].GetSnapshot()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack

	app := tview.NewApplication()

	menu := tview.NewList().
		AddItem("Overview", "", '1', func() {
			showOverview = true
		}).
		AddItem("GDAX", "", '2', func() {
			showOverview = false
		}).
		AddItem("Gemini", "", '3', func() {
			showOverview = false
		}).
		AddItem("Bitfinex", "", '4', func() {
			showOverview = false
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	overview := overview.Table(exchanges)
	// gdaxTable := exchanges["gdax"].Table()

	body := tview.NewFlex().
		SetFullScreen(true).
		AddItem(menu, 20, 1, true).
		AddItem(overview, 0, 1, false)
		// AddItem(gdaxTable, 0, 1, false)

	data := make(chan cq.Quoter, 200)

	go func() {
		exchanges["gdax"].Stream(data)

		for {
			if showOverview == true {

				upd := <-data
				upd.UpdOverviewRow(overview)
				app.Draw()

				time.Sleep(100 * time.Millisecond)
				upd.ClrOverviewBold(overview)
				app.Draw()
			}
		}
	}()

	if err := app.SetRoot(body, true).Run(); err != nil {
		panic(err)
	}
}
