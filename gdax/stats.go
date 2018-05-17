package gdax

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

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
	q := m.data[pair]
	q.ID = pair
	q.High = stats.High
	q.Low = stats.Low
	q.Open = stats.Open
	q.Volume = stats.Volume
	// q.Delta = calcDelta(q.Price, q.Open)
	m.data[pair] = q
	m.Unlock()
}
