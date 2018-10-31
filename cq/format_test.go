package cq

import (
	"testing"

	"github.com/gdamore/tcell"
)

func Test_FmtPair(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"BTC-USD", "BTC/USD"},
		{"ETH-EUR", "ETH/EUR"},
	}

	for _, c := range tc {
		actual := FmtPair(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}

func Test_FmtPrice(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"2300.12345", "   2300.12"},
		{"600.123456", "    600.12"},
		{"9.123456789", "   9.12346"},
		{"9.123", "   9.12300"},
	}

	for _, c := range tc {
		actual := FmtPrice(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}

func Test_FmtDelta(t *testing.T) {
	tc := []struct {
		inprice  string
		inopen   string
		expdelta string
		expcolor tcell.Color
	}{
		{"1000.00", "1200.00", "  -16.67%", tcell.ColorRed},
		{"1000.00", "900.00", "   11.11%", tcell.ColorGreen},
	}

	for _, c := range tc {
		delta, color := FmtDelta(c.inprice, c.inopen)
		if delta != c.expdelta || color != c.expcolor {
			t.Errorf("expected %v and %v; got %v and %v", c.expdelta, c.expcolor, delta, color)
		}
	}
}

func Test_FmtSize(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"12.1234567899", "  12.12345679"},
		{"12.12345", "  12.12345000"},
		{"0.123456789", "   0.12345679"},
		{"0.12345", "   0.12345000"},
	}

	for _, c := range tc {
		actual := FmtSize(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}

func Test_FmtVolume(t *testing.T) {
	tc := []struct {
		input    string
		expected string
	}{
		{"453.12345", "      453"},
		{"453.987654", "      454"},
	}
	for _, c := range tc {
		actual := FmtVolume(c.input)
		if actual != c.expected {
			t.Errorf("expected %v; got %v", c.expected, actual)
		}
	}
}
