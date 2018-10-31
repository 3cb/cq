package bitfinex

import (
	"testing"
)

// func Test_getTickers(t *testing.T) {
// 	var wg *sync.WaitGroup
// 	errCh := make(chan error)

// 	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 	}))
// 	defer srv.Close()

// }

func Test_formatID(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"tBTCUSD", "BTC-USD"},
		{"tETHUSD", "ETH-USD"},
	}

	for _, c := range tc {
		actual := formatID(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}
