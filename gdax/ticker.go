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
