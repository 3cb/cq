package bitfinex

// Quote represents data for a crypto trading pair on bitfinex exchange
// data comes from websocket messages in array format:
//
type Quote struct {
	ID     string
	Price  float32
	Change float32
	Size   float32

	Bid float32
	Ask float32

	Low    float32
	High   float32
	Open   float32
	Volume float32
}
