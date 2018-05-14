package gdax

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

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
	q := m.data[pair]
	q.ID = pair
	q.Price = slice[0].Price
	q.Size = slice[0].Size
	q.Delta = calcDelta(q.Price, q.Open)
	m.data[pair] = q
	m.Unlock()
}
