package coinbase

// Data contains most recent data for each crypto currency pair
// Trailing comments denote which http request or websocket stream the data comes from
// getTrades: https://docs.pro.coinbase.com/#get-trades
// match: https://docs.pro.coinbase.com/#the-matches-channel
// ticker: https://docs.pro.coinbase.com/#the-ticker-channel
// *** GDAX API documentation for websocket ticker channel does not show all available fields as of 2/11/2018
type Data struct {
	Type string `json:"type"` // used to filter websocket messages

	Pair  string `json:"product_id"`
	Price string `json:"price"` // getTrades/match
	Size  string `json:"size"`  // getTrades/match

	Bid string `json:"best_bid"` // getTicker/ticker
	Ask string `json:"best_ask"` // getTicker/ticker

	High   string `json:"high_24h"`   // getStats/ticker
	Low    string `json:"low_24h"`    // getStats/ticker
	Open   string `json:"open_24h"`   // getStats/ticker
	Volume string `json:"volume_24h"` // getStats/ticker
}
