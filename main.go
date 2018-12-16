package main

import (
	"flag"

	"github.com/3cb/cq/bitfinex"
	"github.com/3cb/cq/coinbase"
	"github.com/3cb/cq/cq"
	"github.com/3cb/cq/hitbtc"
	"github.com/3cb/cq/overview"
	"github.com/rivo/tview"
)

func main() {
	theme := flag.String("t", "dark", "color scheme")
	flashCell := flag.Bool("f", true, "flash table cell on trade")
	flag.Parse()

	cfg := cq.SetConfig(theme, flashCell)

	// Initialize exchanges
	exchanges := make(map[string]cq.Exchange)
	exchanges["coinbase"] = coinbase.Init()
	exchanges["bitfinex"] = bitfinex.Init()
	exchanges["hitbtc"] = hitbtc.Init()

	// Create tables with initial data from http requests
	println("Building tables...")
	coinbaseCh := make(chan cq.Quoter, 1)
	bitfinexCh := make(chan cq.Quoter, 1)
	hitbtcCh := make(chan cq.Quoter, 1)

	tables := make(map[string]cq.Quoter)
	go func() {
		coinbaseCh <- exchanges["coinbase"].NewTable(cfg)
	}()
	go func() {
		bitfinexCh <- exchanges["bitfinex"].NewTable(cfg)
	}()
	go func() {
		hitbtcCh <- exchanges["hitbtc"].NewTable(cfg)
	}()

	tables["bitfinex"] = <-bitfinexCh
	tables["hitbtc"] = <-hitbtcCh
	tables["coinbase"] = <-coinbaseCh

	// create overview table and prime with price data
	tables["overview"] = overview.NewTable(cfg)
	for _, ex := range exchanges {
		quotes := ex.GetQuotes()
		for _, q := range quotes {
			msg := cq.UpdateMsg{
				Quote:   q,
				IsTrade: true,
				Flash:   false,
			}
			tables["overview"].InsertQuote(msg, cfg)
		}
	}

	mktView := tables["overview"]

	app := tview.NewApplication()

	view := make(chan cq.Quoter)
	done := make(chan struct{})

	menu := tview.NewList().
		AddItem("Overview", "", '1', func() {
			view <- tables["overview"]
			<-done
		}).
		AddItem("Coinbase", "", '2', func() {
			view <- tables["coinbase"]
			<-done
		}).
		AddItem("Bitfinex", "", '3', func() {
			view <- tables["bitfinex"]
			<-done
		}).
		AddItem("HitBTC", "", '4', func() {
			view <- tables["hitbtc"]
			<-done
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	body := tview.NewFlex().
		SetFullScreen(true).
		AddItem(menu, 20, 1, true).
		AddItem(tables["overview"], 0, 1, false)

	// set colors for menu
	menu.SetBackgroundColor(cfg.Theme.BackgroundColor)
	menu.SetMainTextColor(cfg.Theme.TextColor)
	menu.SetSecondaryTextColor(cfg.Theme.TextColorMenuSecondary)
	menu.SetShortcutColor(cfg.Theme.TextColorMenuShortcut)
	menu.SetSelectedTextColor(cfg.Theme.TextColorMenuSelected)
	menu.SetSelectedBackgroundColor(cfg.Theme.BackgroundColorMenuSelected)

	updateCh, timerCh := cq.StartTimerGroup(exchanges)

	go func() {
		for {
			select {
			case upd := <-updateCh:
				if upd.IsTrade {
					app.QueueUpdateDraw(func() {
						tables[upd.Quote.MarketID].InsertQuote(upd, cfg)
					}).
						QueueUpdateDraw(func() {
							tables["overview"].InsertQuote(upd, cfg)
						})
				} else {
					app.QueueUpdateDraw(func() {
						tables[upd.Quote.MarketID].InsertQuote(upd, cfg)
					})
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
	go func() {
		exchanges["coinbase"].Stream(timerCh)
		exchanges["bitfinex"].Stream(timerCh)
		exchanges["hitbtc"].Stream(timerCh)
		// handle errors here *******************************
	}()

	println("Launching app...")
	if err := app.SetRoot(body, true).Run(); err != nil {
		panic(err)
	}
}
