package gdax

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

// Ticker contains snapshot data from REST API:
// https://docs.gdax.com/#get-product-ticker
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

	api := "https://api.gdax.com/products/" + pair + "/ticker"
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
	q := (m.data[pair]).(Quote)
	q.ID = pair
	q.Bid = ticker.Bid
	q.Ask = ticker.Ask
	m.data[pair] = q
	m.Unlock()
}

// Stats contains 24 hour data from REST API:
// https://docs.gdax.com/#get-24hr-stats
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

	api := "https://api.gdax.com/products/" + pair + "/stats"
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
	q := (m.data[pair]).(Quote)
	q.ID = pair
	q.High = stats.High
	q.Low = stats.Low
	q.Open = stats.Open
	q.Volume = stats.Volume
	m.data[pair] = q
	m.Unlock()
}

func getTrades(m *Market, pair string, wg *sync.WaitGroup, errCh chan error) {
	defer func() {
		wg.Done()
	}()

	slice := []Quote{}

	api := "https://api.gdax.com/products/" + pair + "/trades?limit=1"
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
	q := (m.data[pair]).(Quote)
	q.ID = pair
	q.Price = slice[0].Price
	q.Size = slice[0].Size
	m.data[pair] = q
	m.Unlock()
}
