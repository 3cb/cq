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

	overviewTbl := overview.Table(exchanges)
	gdaxTbl := exchanges["gdax"].Table()

	// showOverview := true
	mktView := overviewTbl

	// handle error slice here
	exchanges["gdax"].GetSnapshot()

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack

	app := tview.NewApplication()

	body := tview.NewFlex().
		SetFullScreen(true)

	menu := tview.NewList().
		AddItem("Overview", "", '1', func() {
			mktView = setView(body, mktView, overviewTbl)
			app.Draw()
		}).
		AddItem("GDAX", "", '2', func() {
			mktView = setView(body, mktView, gdaxTbl)
			app.Draw()
		}).
		AddItem("Gemini", "", '3', func() {
			mktView = setView(body, mktView, overviewTbl)
			app.Draw()
		}).
		AddItem("Bitfinex", "", '4', func() {
			mktView = setView(body, mktView, overviewTbl)
			app.Draw()
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	body.
		AddItem(menu, 20, 1, true).
		AddItem(overviewTbl, 0, 1, false)

	data := make(chan cq.Quoter, 200)

	go func() {
		exchanges["gdax"].Stream(data)

		for {
			if mktView == overviewTbl {
				upd := <-data
				upd.UpdOverviewRow(overviewTbl)
				app.Draw()

				time.Sleep(65 * time.Millisecond)
				upd.ClrOverviewBold(overviewTbl)
				app.Draw()
			}
		}
	}()

	if err := app.SetRoot(body, true).Run(); err != nil {
		panic(err)
	}
}

func setView(body *tview.Flex, mktView *tview.Table, targetView *tview.Table) *tview.Table {
	if mktView != targetView {
		body.RemoveItem(mktView)
		body.AddItem(targetView, 0, 1, false)
	}
	return targetView
}
