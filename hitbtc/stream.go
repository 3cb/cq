package hitbtc

import (
	"errors"
	"strings"

	"github.com/3cb/cq/cq"
	"github.com/gorilla/websocket"
)

// SubscribeMsg contains info to subscribe to websocket data
type SubscribeMsg struct {
	Method string            `json:"method"`
	Params map[string]string `json:"params"`
	ID     string            `json:"id"`
}

// WSMsg contains data from websocket messages
// Params has to be type asserted:
// ticker:
// "params": {
//     "ask": "0.054464",
//     "bid": "0.054463",
//     "last": "0.054463",
//     "open": "0.057133",
//     "low": "0.053615",
//     "high": "0.057559",
//     "volume": "33068.346",
//     "volumeQuote": "1832.687530809",
//     "timestamp": "2017-10-19T15:45:44.941Z",
//     "symbol": "ETHBTC"
//   }
//
// trades:
// "params": {
//     "data": [
//       {
//         "id": 54469813,
//         "price": "0.054670",
//         "quantity": "0.183",
//         "side": "buy",
//         "timestamp": "2017-10-19T16:34:25.041Z"
//       }
//     ],
//     "symbol": "ETHBTC"
//   }
type WSMsg struct {
	VersionJSON string      `json:"jsonrpc"`
	Method      string      `json:"method"`
	Params      interface{} `json:"params"`
}

func connectWS(m *Market, dataCh chan<- cq.Quoter) error {
	var failedSubs []string

	api := "wss://api.hitbtc.com/api/2/ws"

	conn, resp, err := websocket.DefaultDialer.Dial(api, nil)
	if resp.StatusCode != 101 || err != nil {
		return errors.New("unable to connect to hitbtc websocket api")
	}

	m.RLock()
	symbols := m.symbols
	m.RUnlock()

	for _, s := range symbols {
		subTicker := &SubscribeMsg{
			Method: "subscribeTicker",
			Params: map[string]string{
				"symbol": s,
			},
			ID: s,
		}
		subTrades := &SubscribeMsg{
			Method: "subscribeTrades",
			Params: map[string]string{
				"symbol": s,
			},
			ID: s,
		}

		// write ticker sub to websocket
		err := conn.WriteJSON(subTicker)
		if err != nil {
			failedSubs = append(failedSubs, s)
			continue
		}

		// write trades sub to websocket
		err = conn.WriteJSON(subTrades)
		if err != nil {
			failedSubs = append(failedSubs, s)
		}
	}

	if len(failedSubs) > 0 {
		b := strings.Builder{}
		b.WriteString("failed to subscribe to the following symbols: ")
		for i, fail := range failedSubs {
			b.WriteString(fail)
			if i < len(failedSubs)-1 {
				b.WriteString(", ")
			}
		}
		println(b.String())
		return errors.New(b.String())
	}

	go func() {
		defer func() {
			conn.Close()
			m.streaming = false
		}()

		m.streaming = true

		for {
			var msg WSMsg

			err := conn.ReadJSON(&msg)
			if err != nil {
				return
			}

			switch msg.Method {
			case "ticker":
				p := (msg.Params).(map[string]interface{})
				s := (p["symbol"]).(string)

				m.Lock()
				q := (m.data[s]).(Quote)
				q.Ask = (p["ask"]).(string)
				q.Bid = (p["bid"]).(string)
				q.Low = (p["low"]).(string)
				q.High = (p["high"]).(string)
				q.Open = (p["open"]).(string)
				q.Volume = (p["volume"]).(string)
				m.data[s] = q
				dataCh <- q
				m.Unlock()
			case "updateTrades":
				p := (msg.Params).(map[string]interface{})
				d := (p["data"]).([]interface{})
				u := (d[0]).(map[string]interface{})
				s := (p["symbol"]).(string)

				m.Lock()
				q := (m.data[s]).(Quote)
				q.Price = (u["price"]).(string)
				q.Size = (u["quantity"]).(string)
				m.data[s] = q
				dataCh <- q
				m.Unlock()
			default:
				continue
			}
		}
	}()

	return nil
}
