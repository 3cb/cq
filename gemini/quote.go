package gemini

// Quote contains most recent data for each crypto currency pair
type Quote struct {
	Symbol string

	ID     string
	Price  string
	Bid    string
	Ask    string
	Volume string
}
