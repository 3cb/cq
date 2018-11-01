package bitfinex

import (
	"testing"
)

func Test_fmtID(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"tBTCUSD", "BTC/USD"},
		{"tETHUSD", "ETH/USD"},
	}

	for _, c := range tc {
		actual := fmtID(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}
