package main

import (
	"flag"

	"github.com/3cb/cq/gdax"
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
	exchanges["gdax"].Print()
}

type exchange interface {
	Snapshot() []error
	Stream() error
	Update() error
	Print()
}
