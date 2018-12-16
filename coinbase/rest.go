package coinbase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// getSnapshot method performs concurrent http requests the GDAX REST API to get initial
// market data
func (m *Market) getSnapshot() []error {
	var e []error

	m.RLock()
	pairs := m.pairs
	m.RUnlock()
	l := len(pairs)
	errCh := make(chan error, (9 * l))

	wg := &sync.WaitGroup{}
	wg.Add(3 * l)
	// break requests up into bursts to satisfy coinbase throttling
	for i := 0; i < 5; i++ {
		go getTrades(m, pairs[i], wg, errCh)
		go getStats(m, pairs[i], wg, errCh)
		go getTicker(m, pairs[i], wg, errCh)
	}
	time.Sleep(1 * time.Second)
	for i := 5; i < 10; i++ {
		go getTrades(m, pairs[i], wg, errCh)
		go getStats(m, pairs[i], wg, errCh)
		go getTicker(m, pairs[i], wg, errCh)
	}
	time.Sleep(1 * time.Second)
	for i := 10; i < 15; i++ {
		go getTrades(m, pairs[i], wg, errCh)
		go getStats(m, pairs[i], wg, errCh)
		go getTicker(m, pairs[i], wg, errCh)
	}
	wg.Wait()

	close(errCh)
	for err := range errCh {
		e = append(e, err)
	}

	return e
}

// Ticker contains snapshot data from REST API:
// https://docs.pro.coinbase.com/#get-product-ticker
type Ticker struct {
	ID  string
	Bid string `json:"bid"`
	Ask string `json:"ask"`
}

func getTicker(m *Market, pair string, wg *sync.WaitGroup, errCh chan error) {
	defer func() {
		wg.Done()
	}()

	ticker := Ticker{}

	api := "https://api.pro.coinbase.com/products/" + pair + "/ticker"
	resp, err := http.Get(api)
	if err != nil {
		errCh <- err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
	}
	err = json.Unmarshal(bytes, &ticker)
	if err != nil {
		errCh <- err
	}

	m.Lock()
	q := m.data[pair]
	q.Bid = ticker.Bid
	q.Ask = ticker.Ask
	m.data[pair] = q
	m.Unlock()
}

// Stats contains 24 hour data from REST API:
// https://docs.pro.coinbase.com/#get-24hr-stats
type Stats struct {
	ID     string
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Volume string `json:"volume"`
}

func getStats(m *Market, pair string, wg *sync.WaitGroup, errCh chan error) {
	defer func() {
		wg.Done()
	}()

	stats := Stats{}

	api := "https://api.pro.coinbase.com/products/" + pair + "/stats"
	resp, err := http.Get(api)
	if err != nil {
		errCh <- err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errCh <- err
	}
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		errCh <- err
	}

	m.Lock()
	q := m.data[pair]
	q.High = stats.High
	q.Low = stats.Low
	q.Open = stats.Open
	q.Volume = stats.Volume
	m.data[pair] = q
	m.Unlock()
}

// https://docs.pro.coinbase.com/#get-trades
func getTrades(m *Market, pair string, wg *sync.WaitGroup, errCh chan error) {
	defer func() {
		wg.Done()
	}()

	slice := []Data{}

	api := "https://api.pro.coinbase.com/products/" + pair + "/trades?limit=1"
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
	err = json.Unmarshal(bytes, &slice)
	if err != nil {
		errCh <- err
		return
	}

	m.Lock()
	q := m.data[pair]
	q.Price = slice[0].Price
	q.Size = slice[0].Size
	m.data[pair] = q
	m.Unlock()
}
