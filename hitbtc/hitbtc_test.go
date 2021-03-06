package hitbtc

import (
	"testing"
)

func Test_fmtID(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"BTCUSD", "BTC/USD"},
		{"BCHUSD", "BCH/USD"},
	}

	for _, c := range tc {
		if actual := fmtID(c.input); actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}

	}
}
