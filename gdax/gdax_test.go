package gdax

import "testing"

func Test_fmtID(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"BTC-USD", "BTC/USD"},
		{"ETH-EUR", "ETH/EUR"},
	}

	for _, c := range tc {
		actual := fmtID(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}
