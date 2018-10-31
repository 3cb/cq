package gemini

import "testing"

func Test_getVolDenom(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"btcusd", "BTC"},
		{"ethusd", "ETH"},
		{"ltcbtc", "LTC"},
	}

	for _, c := range tc {
		actual := getVolDenom(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}

func Test_setID(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"btcusd", "BTC-USD"},
		{"ethusd", "ETH-USD"},
	}

	for _, c := range tc {
		actual := setID(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}
