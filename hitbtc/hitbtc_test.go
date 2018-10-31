package hitbtc

import (
	"testing"
)

func Test_formatID(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"BTCUSD", "BTC-USD"},
		{"BCHUSD", "BCH-USD"},
	}

	for _, c := range tc {
		if actual := formatID(c.input); actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}

	}
}
