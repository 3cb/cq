package cq

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

// All formatting functions pad the front of strings with empty spaces
// to maintain column width.  Fixed column widths are set according to value type:
// Price: 10
// Delta: 9
// Size: 13
// Volume: 9

// FmtPrice formats price data for display
// If price is >= 10 it uses 2 decimal places
// If price is below 10 it uses 5 decimal places
func FmtPrice(price string) string {
	num, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "no data"
	}
	if num >= 10 {
		num = float64(int64(num*100+0.5)) / 100
		price = fmt.Sprintf("%.2f", num)
	} else {
		num = float64(int64(num*100000+0.5)) / 100000
		price = fmt.Sprintf("%.5f", num)
	}
	spc := 10 - len(price)
	b := strings.Builder{}
	for i := 0; i < spc; i++ {
		b.WriteString(" ")
	}
	b.WriteString(price)
	return b.String()
}

// FmtDelta calculates price delta and provides appropriate formatting
func FmtDelta(price string, open string, cfg Config) (string, tcell.Color) {
	if len(price) > 0 && len(open) > 0 {
		p, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return "no data", cfg.Theme.TextColor
		}
		o, err := strconv.ParseFloat(open, 64)
		if err != nil {
			return "no data", cfg.Theme.TextColor

		}
		d := (p - o) / o * 100
		delta := strconv.FormatFloat(d, 'f', 2, 64)
		spc := 9 - len(delta) - 1 // extra "1" acounts for "%" character
		b := strings.Builder{}
		for i := 0; i < spc; i++ {
			b.WriteString(" ")
		}
		b.WriteString(delta)
		b.WriteString("%")
		if d >= 0 {
			return b.String(), cfg.Theme.TextColorUp
		}
		return b.String(), cfg.Theme.TextColorDown
	}
	return "no data", cfg.Theme.TextColor
}

// FmtSize formats trade size data with 8 decimal places
func FmtSize(size string) string {
	num, err := strconv.ParseFloat(size, 64)
	if err != nil {
		return "no data"
	}
	num = float64(int64(num*100000000+0.5)) / 100000000
	size = fmt.Sprintf("%.8f", num)
	spc := 13 - len(size)
	b := strings.Builder{}
	for i := 0; i < spc; i++ {
		b.WriteString(" ")
	}
	b.WriteString(size)
	return b.String()
}

// FmtVolume formats volume data by rounding to nearest whole number
func FmtVolume(vol string) string {
	num, err := strconv.ParseFloat(vol, 64)
	if err != nil {
		return "no data"
	}
	vol = fmt.Sprint(int64(num + 0.5))
	spc := 9 - len(vol)
	b := strings.Builder{}
	for i := 0; i < spc; i++ {
		b.WriteString(" ")
	}
	b.WriteString(vol)
	return b.String()
}
