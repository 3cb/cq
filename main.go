package main

import (
	"flag"
)

func main() {
	snap := flag.Bool("snap", false, "get quote snapshot")
	flag.Parse()

	exchanges := make(map[string]*exchange)

	exchanges["gdax"] = &exchange{
		streaming: false,
		pairs:     []string{"BTC-USD", "BTC-EUR", "BTC-GBP", "BCH-USD", "BCH-BTC", "BCH-EUR", "ETH-USD", "ETH-BTC", "ETH-EUR", "LTC-USD", "LTC-BTC", "LTC-EUR"},
	}
	exchanges["gemini"] = &exchange{
		streaming: false,
		pairs:     []string{"btcusd", "ethusd", "ethbtc"},
	}
	exchanges["bitfinex"] = &exchange{
		streaming: false,
		pairs:     []string{"btcusd", "btceur", "btcgbp", "btcjpy", "ethusd", "ethbtc", "etheur", "ethgbp", "ethjpy", "bchusd", "bchbtc", "bcheth", "ltcusd", "ltcbtc"},
	}
}

type exchange struct {
	streaming bool
	pairs     []string
}
