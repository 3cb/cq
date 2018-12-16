package bitfinex

import (
	"testing"

	"github.com/3cb/cq/cq"
)

func Test_streaming(t *testing.T) {
	t.Run("test registerChannelID", func(t *testing.T) {
		tc := []map[string]interface{}{
			{
				"event":   "subscribed",
				"channel": "ticker",
				"chanId":  float64(173),
				"symbol":  "tBTCUSD",
				"pair":    "BTCUSD",
			},
			{
				"event":   "subscribed",
				"channel": "trades",
				"chanId":  float64(371),
				"symbol":  "tBTCUSD",
				"pair":    "BTCUSD",
			},
		}
		store2 := ChannelStore{
			Trades:  make(map[float64]string),
			Tickers: make(map[float64]string),
		}

		registerChannelID(tc[0], &store2)
		symbol, ok := store2.Tickers[173]
		if !ok {
			t.Errorf("expected key(173); does not exist in map.")
		}
		if symbol != "tBTCUSD" {
			t.Errorf("expected tBTCUSD; got %v", symbol)
		}

		registerChannelID(tc[1], &store2)
		sym, ok := store2.Trades[371]
		if !ok {
			t.Errorf("expected key(371); does not exist in map.")
		}
		if sym != "tBTCUSD" {
			t.Errorf("expected tBTCUSD; got %v", sym)
		}
	})

	t.Run("test queryStore", func(t *testing.T) {
		store := &ChannelStore{
			Trades:  make(map[float64]string),
			Tickers: make(map[float64]string),
		}

		// give store some data
		store.Trades[123] = "tBTCUSD"
		store.Trades[222] = "tETHEUR"
		store.Tickers[155] = "tBTCUSD"
		store.Tickers[333] = "tETHEUR"

		tc := []struct {
			input      float64
			expPair    string
			expChannel string
		}{
			{123, "tBTCUSD", "trade"},
			{222, "tETHEUR", "trade"},
			{155, "tBTCUSD", "ticker"},
			{333, "tETHEUR", "ticker"},
		}

		for _, c := range tc {
			pair, channel := queryStore(store, c.input)
			if pair != c.expPair || channel != c.expChannel {
				t.Errorf("expected %v and %v; got %v and %v", c.expPair, c.expChannel, pair, channel)
			}
		}
	})

	t.Run("test tickerToQuote", func(t *testing.T) {
		m := Init()
		dataCh := make(chan cq.TimerMsg)
		pair := "tBTCUSD"
		upd := []interface{}{
			6615.5,
			41.77077901,
			6617.7,
			42.87899145,
			-112.7,
			-0.0167,
			6618.1,
			11108.82107481,
			6779.9,
			6609.3,
		}
		expected := cq.Quote{
			Bid:        "6615.5",
			Ask:        "6617.7",
			Change:     "-112.7",
			ChangePerc: "-0.0167",
			Volume:     "11108.82107481",
			High:       "6779.9",
			Low:        "6609.3",
		}

		go tickerToQuote(m, upd, pair, dataCh)
		actualCh := <-dataCh
		actual := m.data[pair]

		if actual.Bid != expected.Bid || actualCh.Bid != expected.Bid {
			t.Errorf("expected %v and %v; got %v and %v", expected.Bid, expected.Bid, actual.Bid, actualCh.Bid)
		}
		if actual.Ask != expected.Ask || actualCh.Ask != expected.Ask {
			t.Errorf("expected %v and %v; got %v and %v", expected.Ask, expected.Ask, actual.Ask, actualCh.Ask)
		}
		if actual.Change != expected.Change || actualCh.Change != expected.Change {
			t.Errorf("expected %v and %v; got %v and %v", expected.Change, expected.Change, actual.Change, actualCh.Change)
		}
		if actual.ChangePerc != expected.ChangePerc || actualCh.ChangePerc != expected.ChangePerc {
			t.Errorf("expected %v and %v; got %v and %v", expected.ChangePerc, expected.ChangePerc, actual.ChangePerc, actualCh.ChangePerc)
		}
		if actual.Volume != expected.Volume || actualCh.Volume != expected.Volume {
			t.Errorf("expected %v and %v; got %v and %v", expected.Volume, expected.Volume, actual.Volume, actualCh.Volume)
		}
		if actual.High != expected.High || actualCh.High != expected.High {
			t.Errorf("expected %v and %v; got %v and %v", expected.High, expected.High, actual.High, actualCh.High)
		}
		if actual.Low != expected.Low || actualCh.Low != expected.Low {
			t.Errorf("expected %v and %v; got %v and %v", expected.Low, expected.Low, actual.Low, actualCh.Low)
		}
	})

	t.Run("test tradeToQuote", func(t *testing.T) {
		m := Init()
		dataCh2 := make(chan cq.TimerMsg)
		pair := "tBTCUSD"
		tc := [][]interface{}{
			{306300718, 1540146321621, -0.01997, 6625.6},
			{306300697, 1540146302057, 0.10700401, 6625.7},
			{306300690, 1540146285976, 0.015513, 6625.7},
			{306300688, 1540146282272, -0.05, 6625.6},
		}
		expected := []cq.Quote{
			{
				Price: "6625.6",
				Size:  "0.01997",
			},
			{
				Price: "6625.7",
				Size:  "0.10700401",
			},
			{
				Price: "6625.7",
				Size:  "0.015513",
			},
			{
				Price: "6625.6",
				Size:  "0.05",
			},
		}

		for i, c := range tc {
			go tradeToQuote(m, c, pair, dataCh2)
			actualCh := <-dataCh2
			actual := m.data[pair]

			if actual.Price != expected[i].Price || actualCh.Price != expected[i].Price {
				t.Errorf("expected %v and %v; got %v and %v", expected[i].Price, expected[i].Price, actual.Price, actualCh.Price)
			}
			if actual.Size != expected[i].Size || actualCh.Size != expected[i].Size {
				t.Errorf("expected %v and %v; got %v and %v", expected[i].Size, expected[i].Size, actual.Size, actualCh.Size)
			}
		}
	})

}
