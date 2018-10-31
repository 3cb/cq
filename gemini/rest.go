package gemini

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// TickerResponse holds data from ticker GET request
type TickerResponse struct {
	Ask    string                 `json:"ask"`
	Bid    string                 `json:"bid"`
	Last   string                 `json:"last"`
	Volume map[string]interface{} `json:"volume"`
}

// GET https://api.gemini.com/v1/pubticker/:symbol
// {
//     "ask": "977.59",
//     "bid": "977.35",
//     "last": "977.65",
//     "volume": {
//         "BTC": "2210.505328803",
//         "USD": "2135477.463379586263",
//         "timestamp": 1483018200000
//     }
// }
func (m *Market) getTicker(pair string, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	var data TickerResponse

	api := "https://api.gemini.com/v1/pubticker/" + pair
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
	q := (m.data[pair]).(Quote)
	q.Symbol = pair
	q.ID = setID(pair)
	q.Bid = data.Bid
	q.Ask = data.Ask
	q.Volume = (data.Volume[getVolDenom(pair)]).(string)
	m.data[pair] = q
	m.Unlock()
}

func setID(pair string) string {
	temp := strings.Split(pair, "")
	temp2 := []string{temp[0], temp[1], temp[2], "-", temp[3], temp[4], temp[5]}
	return strings.ToUpper(strings.Join(temp2, ""))
}

// getVolDenom takes the trading pair as a parameter and returns the denomination
// for volume data
func getVolDenom(pair string) string {
	var denom []string

	temp := strings.Split(pair, "")
	denom = append(denom, temp[:3]...)
	return strings.ToUpper(strings.Join(denom, ""))
}

// Trade contains data from Gemini trades rest request
type Trade struct {
	Timestamp   int    `json:"timestamp"`
	TimestampMS int    `json:"timestampms"`
	TID         int    `json:"tid"`
	Price       string `json:"price"`
	Amount      string `json:"amount"`
	Exchange    string `json:"exchange"`
	Type        string `json:"type"`
}

// GET https://api.gemini.com/v1/trades/:symbol
// [
//   {
//     "timestamp": 1420088400,
//     "timestampms": 1420088400122,
//     "tid": 155814,
//     "price": "822.12",
//     "amount": "12.10",
//     "exchange": "gemini",
//     "type": "buy"
//   },
//   ...
// ]
func (m *Market) getTrades(pair string, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()
	var data []Trade

	api := "https://api.gemini.com/v1/trades/" + pair + "?limit_trades=1"

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

	if len(data) > 0 {
		lastTrade := data[0]

		m.Lock()
		q := (m.data[pair]).(Quote)
		q.Price = lastTrade.Price
		q.Size = lastTrade.Amount
		m.data[pair] = q
		m.Unlock()
	} else {
		errCh <- fmt.Errorf("no data from gemini trades request for %v", pair)
	}
}
