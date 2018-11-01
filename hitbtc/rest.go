package hitbtc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

// TickerEntry holds data for element of ticker response array
// https://api.hitbtc.com/api/2/public/ticker
type TickerEntry struct {
	Symbol      string `json:"symbol"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Low         string `json:"low"`
	High        string `json:"high"`
	Open        string `json:"open"`
	Volume      string `json:"volume"`
	VolumeQuote string `json:"volumeQuote"`
	Timestamp   string `json:"timestamp"`
}

// TradesEntry holds data for element of trades response array
// https://api.hitbtc.com/api/2/public/trades/{symbol}
type TradesEntry struct {
	ID        int    `json:"id"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Side      string `json:"side"`
	Timestamp string `json:"timestamp"`
}

// get initial data for overview and hitbtc tables
func (m *Market) getSnapshot() []error {
	var e []error

	m.RLock()
	symbols := m.symbols
	m.RUnlock()
	l := len(symbols)
	errCh := make(chan error, l+1)

	wg := &sync.WaitGroup{}
	wg.Add(l + 1)
	go getTicker(m, wg, errCh)
	for _, s := range symbols {
		go getTrades(m, s, wg, errCh)
	}
	wg.Wait()

	close(errCh)
	for err := range errCh {
		e = append(e, err)
	}

	return e
}

// get most recent ticker values for all symbols:
// https://api.hitbtc.com/api/2/public/ticker
func getTicker(m *Market, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	var data []TickerEntry

	api := "https://api.hitbtc.com/api/2/public/ticker"

	resp, err := http.Get(api)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		errCh <- err
		return
	}

	for _, s := range m.symbols {
		for _, d := range data {
			if s == d.Symbol {
				m.Lock()
				q := (m.data[s]).(Quote)
				q.Symbol = s
				q.Ask = d.Ask
				q.Bid = d.Bid
				q.Low = d.Low
				q.High = d.High
				q.Open = d.Open
				q.Volume = d.Volume
				m.data[s] = q
				m.Unlock()
			}
		}
	}
}

// get last trade price and volume data by symbol:
// https://api.hitbtc.com/api/2/public/trades/{symbol}
func getTrades(m *Market, symbol string, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	var data []TradesEntry

	// get only most recent trade
	api := "https://api.hitbtc.com/api/2/public/trades/" + symbol + "?limit=1"

	resp, err := http.Get(api)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		errCh <- err
		return
	}

	m.Lock()
	q := (m.data[symbol]).(Quote)
	q.Price = data[0].Price
	q.Size = data[0].Quantity
	m.data[symbol] = q
	m.Unlock()
}
