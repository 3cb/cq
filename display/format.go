package display

import (
	"fmt"
	"strconv"
	"strings"
)

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
		return fmt.Sprintf("%.2f", num)

	}
	num = float64(int64(num*100000+0.5)) / 100000
	return fmt.Sprintf("%.5f", num)
}

// FmtDelta calculates price delta and provides appropriate formatting
func FmtDelta(price string, open string) string {
	if len(price) > 0 && len(open) > 0 {
		p, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return "no data"
		}
		o, err := strconv.ParseFloat(open, 64)
		if err != nil {
			return "no data"
		}
		d := (p - o) / o * 100
		b := strings.Builder{}
		b.WriteString(strconv.FormatFloat(d, 'f', 2, 64))
		b.WriteString("%")
		return b.String()
	}
	return "no data"
}

// FmtSize formats trade size data with 8 decimal places
func FmtSize(size string) string {
	num, err := strconv.ParseFloat(size, 64)
	if err != nil {
		return "no data"
	}
	num = float64(int64(num*100000000+0.5)) / 100000000
	return fmt.Sprintf("%.8f", num)

}

// FmtVolume formats volume data by rounding to nearest whole number
func FmtVolume(vol string) string {
	num, err := strconv.ParseFloat(vol, 64)
	if err != nil {
		return "no data"
	}
	return fmt.Sprint(int64(num + 0.5))
}
