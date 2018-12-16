package cq

// Quote contains all quote data as strings to support display
type Quote struct {
	MarketID   string
	ID         string
	Price      string
	Change     string
	ChangePerc string
	Size       string
	Bid        string
	Ask        string
	Low        string
	High       string
	Open       string
	Volume     string
}

// UpdateMsg carries quotes from TimerGroup event loop to main cq event loop
// IsTrade and Flash fields allow event loop to set table fonts for quotes
type UpdateMsg struct {
	Quote

	// true for "trade" false for "ticker"
	IsTrade bool

	// true if entire cell should flash on new trade
	// false if only text should flash
	Flash bool
}

// TimerMsg is used to carry quotes from package websocket loops
// to TimerGroup event loop
type TimerMsg struct {
	Quote

	// true for "trade" false for "ticker"
	IsTrade bool
}
