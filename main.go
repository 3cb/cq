package main

import (
	"github.com/3cb/cq/bitfinex"
	"github.com/3cb/cq/coinbase"
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/gemini"
	"github.com/3cb/cq/hitbtc"
	"github.com/3cb/cq/overview"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	// Initialize exchanges
	exchanges := make(map[string]cq.Exchange)
	exchanges["coinbase"] = coinbase.Init()
	exchanges["bitfinex"] = bitfinex.Init()
	exchanges["hitbtc"] = hitbtc.Init()
	exchanges["gemini"] = gemini.Init()

	// Create tables with initial data from http requests
	println("Building tables...")
	coinbaseCh := make(chan *tview.Table)
	bitfinexCh := make(chan *tview.Table)
	hitbtcCh := make(chan *tview.Table)
	geminiCh := make(chan *tview.Table)

	overviewTbl := overview.Table()
	go func() {
		coinbaseCh <- exchanges["coinbase"].Table(overviewTbl)
	}()
	go func() {
		bitfinexCh <- exchanges["bitfinex"].Table(overviewTbl)
	}()
	go func() {
		hitbtcCh <- exchanges["hitbtc"].Table(overviewTbl)
	}()
	go func() {
		geminiCh <- exchanges["gemini"].Table(overviewTbl)
	}()
	coinbaseTbl := <-coinbaseCh
	bitfinexTbl := <-bitfinexCh
	hitbtcTbl := <-hitbtcCh
	geminiTbl := <-geminiCh

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
		AddItem("Coinbase", "", '2', func() {
			view <- coinbaseTbl
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

	updateCh, routerCh := cq.NewTimerGroup(exchanges)

	go func() {
		for {
			select {
			case upd := <-updateCh:
				var t *tview.Table
				switch upd.Quote.MarketID() {
				case "coinbase":
					t = coinbaseTbl
				case "bitfinex":
					t = bitfinexTbl
				case "hitbtc":
					t = hitbtcTbl
				case "gemini":
					t = geminiTbl
				}

				switch upd.UpdType {
				case "ticker":
					// println("ticker")
					app.QueueUpdate(upd.Quote.TickerUpdate(t))
					app.Draw()
				case "trade":
					// println("trade")
					if upd.Flash == true {
						app.QueueUpdate(upd.Quote.TradeUpdate(overviewTbl, t, tcell.AttrBold))
					} else {
						app.QueueUpdate(upd.Quote.TradeUpdate(overviewTbl, t, tcell.AttrNone))
					}
					app.Draw()
				}

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
	exchanges["coinbase"].Stream(routerCh)
	exchanges["bitfinex"].Stream(routerCh)
	exchanges["hitbtc"].Stream(routerCh)
	exchanges["gemini"].Stream(routerCh)
	// handle errors here *******************************

	println("Launching app...")
	if err := app.SetRoot(body, true).Run(); err != nil {
		panic(err)
	}
}
