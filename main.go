package main

import (
	"time"

	"github.com/3cb/cq/bitfinex"
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/gdax"
	"github.com/3cb/cq/gemini"
	"github.com/3cb/cq/hitbtc"
	"github.com/3cb/cq/overview"
	"github.com/3cb/tview"
	"github.com/gdamore/tcell"
)

func main() {
	// Initialize exchanges
	exchanges := make(map[string]cq.Exchange)
	exchanges["gdax"] = gdax.Init()
	exchanges["bitfinex"] = bitfinex.Init()
	exchanges["hitbtc"] = hitbtc.Init()
	exchanges["gemini"] = gemini.Init()

	// Create tables with initial data from http requests
	println("Building tables...")
	overviewTbl := overview.Table()
	gdaxTbl := exchanges["gdax"].Table(overviewTbl)
	bitfinexTbl := exchanges["bitfinex"].Table(overviewTbl)
	hitbtcTbl := exchanges["hitbtc"].Table(overviewTbl)
	geminiTbl := exchanges["gemini"].Table(overviewTbl)
	mktView := overviewTbl

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack

	app := tview.NewApplication()

	view := make(chan *tview.Table)
	done := make(chan struct{})

	menu := tview.NewList().
		AddItem("Overview", "", '1', func() {
			view <- overviewTbl
			<-done
		}).
		AddItem("GDAX", "", '2', func() {
			view <- gdaxTbl
			<-done
		}).
		AddItem("Bitfinex", "", '3', func() {
			view <- bitfinexTbl
			<-done
		}).
		AddItem("BitBTC", "", '4', func() {
			view <- hitbtcTbl
			<-done
		}).
		AddItem("Gemini", "", '5', func() {
			view <- geminiTbl
			<-done
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	body := tview.NewFlex().
		SetFullScreen(true).
		AddItem(menu, 20, 1, true).
		AddItem(overviewTbl, 0, 1, false)

	data := make(chan cq.Quoter, 500)

	go func() {
		for {
			select {
			case upd := <-data:
				var t *tview.Table
				switch upd.MarketID() {
				case "gdax":
					t = gdaxTbl
				case "bitfinex":
					t = bitfinexTbl
				case "hitbtc":
					t = hitbtcTbl
				case "gemini":
					t = geminiTbl
				}
				app.QueueUpdate(upd.UpdOverviewRow(overviewTbl))
				app.QueueUpdate(upd.UpdRow(t))

				time.Sleep(85 * time.Millisecond)
				app.QueueUpdate(upd.ClrOverviewBold(overviewTbl))
				app.QueueUpdate(upd.ClrBold(t))
			case tbl := <-view:
				if mktView != tbl {
					body.RemoveItem(mktView)
					body.AddItem(tbl, 0, 1, false)
					mktView = tbl
				}
				done <- struct{}{}
			}
		}
	}()

	println("Connecting to exchanges...")
	// handle errors here *******************************
	exchanges["gdax"].Stream(data)
	exchanges["bitfinex"].Stream(data)
	exchanges["hitbtc"].Stream(data)
	exchanges["gemini"].Stream(data)
	// handle errors here *******************************

	println("Launching app...")
	if err := app.SetRoot(body, true).Run(); err != nil {
		panic(err)
	}
}
