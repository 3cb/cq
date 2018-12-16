package cq

// Exchange interface allows caller to get http snapshot price quotes,
// stream live data via websocket, create tables for display in gui
type Exchange interface {
	// GetIDs returns a slice of trading pairs formatted with "/" (i.e., BTC/USD)
	GetIDs() []string

	// NewTable returns display table with initial data
	NewTable(Config) Quoter

	// GetQuotes returns a map used to prime overview table with data
	// Keys are pair IDs separated with "/". Values are of type Quote.
	GetQuotes() map[string]Quote

	// Stream launches goroutine to stream price data to display table
	Stream(chan<- TimerMsg) error
}
