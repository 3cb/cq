package main

import (
	"flag"

	"github.com/3cb/cq/gdax"
	"github.com/rivo/tview"
)

func main() {
	snap := flag.Bool("snap", false, "get quote snapshot")
	flag.Parse()

	exchanges := make(map[string]exchange)
	exchanges["gdax"] = gdax.Init()
	// exchanges["gemini"] = gemini.Init()
	// exchanges["bitfinex"] = bitfinex.Init()

	if *snap {
		// handle error slice here
		exchanges["gdax"].Snapshot()
	}
	// exchanges["gdax"].Print()

	app := tview.NewApplication()

	list := tview.NewList().
		AddItem("Overview", "", '1', nil).
		AddItem("GDAX", "", '2', nil).
		AddItem("Gemini", "", '3', nil).
		AddItem("Bitfinex", "", '4', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	table := exchanges["gdax"].Table()

	// menu := tview.NewFrame(list).
	// 	SetBorders(0, 0, 0, 0, 0, 0).
	// 	SetBorder(true)

	flex := tview.NewFlex().
		AddItem(list, 30, 1, true).
		AddItem(table, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

type exchange interface {
	Snapshot() []error
	Stream() error
	Table() *tview.Table
	Print()
}
