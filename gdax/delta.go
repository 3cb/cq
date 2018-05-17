package gdax

import "strconv"

func calcDelta(price string, open string) string {
	delta := ""
	if len(price) > 0 && len(open) > 0 {
		p, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return delta
		}
		o, err := strconv.ParseFloat(open, 64)
		if err != nil {
			return delta
		}
		d := (p - o) / o * 100
		delta = strconv.FormatFloat(d, 'E', -1, 64)
	}
	return delta
}
