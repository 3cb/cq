package bitfinex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func (m *Market) getTickers(errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	var data [][]interface{}

	api := "https://api.bitfinex.com/v2/tickers?symbols="
	b := strings.Builder{}
	b.WriteString(api)
	for i, pair := range m.pairs {
		b.WriteString(pair)
		if i < len(m.pairs)-1 {
			b.WriteString(",")
		}
	}
	api = b.String()

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
	for _, v := range data {
		symbol := (v[0]).(string)
		q := m.data[symbol].(Quote)
		q.Symbol = symbol
		q.ID = formatID(symbol)
		q.Bid = (v[1]).(float64)
		q.Ask = (v[3]).(float64)
		q.Change = (v[5]).(float64)
		q.ChangePerc = (v[6]).(float64)
		q.Volume = (v[8]).(float64)
		q.High = (v[9]).(float64)
		q.Low = (v[10]).(float64)
		m.data[symbol] = q
	}
	m.Unlock()
}

func (m *Market) getTrades(pair string, errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	var data [][]interface{}

	api := fmt.Sprintf("https://api.bitfinex.com/v2/trades/%v/hist?limit=1", pair)
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
	for _, val := range data {
		q := (m.data[pair]).(Quote)
		q.Price = (val[3]).(float64)
		q.Size = (val[2]).(float64)
		m.data[pair] = q
	}
	m.Unlock()
}

func formatID(symbol string) string {
	s1 := strings.Split(symbol, "")
	b := strings.Builder{}
	for i := 1; i < 4; i++ {
		b.WriteString(s1[i])
	}
	b.WriteString("-")
	for i := 4; i < 7; i++ {
		b.WriteString(s1[i])
	}
	return b.String()
}
