package bitfinex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func (m *Market) getTickers(errCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	var data [][]interface{}

	api := "https://api.bitfinex.com/v2/tickers?symbols="
	b := strings.Builder{}
	b.WriteString(api)
	m.Lock()
	pairs := m.pairs
	m.Unlock()
	for i, pair := range pairs {
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
		q := m.data[symbol]
		q.Bid = strconv.FormatFloat((v[1]).(float64), 'f', -1, 64)
		q.Ask = strconv.FormatFloat((v[3]).(float64), 'f', -1, 64)
		q.Change = strconv.FormatFloat((v[5]).(float64), 'f', -1, 64)
		q.ChangePerc = strconv.FormatFloat((v[6]).(float64), 'f', -1, 64)
		q.Volume = strconv.FormatFloat((v[8]).(float64), 'f', -1, 64)
		q.High = strconv.FormatFloat((v[9]).(float64), 'f', -1, 64)
		q.Low = strconv.FormatFloat((v[10]).(float64), 'f', -1, 64)
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
		q := m.data[pair]
		q.Price = strconv.FormatFloat((val[3]).(float64), 'f', -1, 64)
		q.Size = strconv.FormatFloat(math.Abs((val[2]).(float64)), 'f', -1, 64)
		m.data[pair] = q
	}
	m.Unlock()
}
