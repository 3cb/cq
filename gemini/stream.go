package gemini

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/3cb/cq/cq"
	"github.com/3cb/ssc"
)

// WSMessage holds data from gemini websocket message
// Only pay attention to trade events:
// "events":[
//    // The first event is the trade event representing
//    // the volume executed in this trade
//       {
//          "type":"trade",
//          "tid":1111597035,
//          "price":"2559.98",
//          "amount":"0.07365713",
//          "makerSide":"ask"
//       },
//       // the second event is a change event which decrements the
//       // quantity of an order that remains on the ask side of the book
//       {
//          "type":"change",
//          "side":"ask",
//          "price":"2559.98",
//          "remaining":"20.98651537",
//          "delta":"-0.07365713",
//          "reason":"trade"
//       }
//    ]
type WSMessage struct {
	Type        string                   `json:"type"`
	Sequence    int                      `json:"socket_sequence"`
	EventID     int                      `json:"eventId"`
	Events      []map[string]interface{} `json:"events"`
	Timestamp   int                      `json:"timestamp"`
	TimestampMS int                      `json:"timestampms"`
}

// wss://api.gemini.com/v1/marketdata/:symbol
func connectWS(m *Market, data chan<- cq.Quoter) error {
	var sockets []string

	m.Lock()
	pairs := m.pairs
	m.Unlock()

	for _, pair := range pairs {
		api := "wss://api.gemini.com/v1/marketdata/" + pair + "?top_of_book=true"
		sockets = append(sockets, api)
	}

	pool := ssc.NewPool(sockets, time.Second*45)
	if err := pool.Start(); err != nil {
		return err
	}

	go func() {
		for {
			var msg WSMessage

			v := <-pool.Outbound
			if v.Type == 1 {
				err := json.Unmarshal(v.Payload, &msg)
				if err != nil {
					return
				}
				urlSlice := strings.Split(v.ID, "/")
				temp := urlSlice[len(urlSlice)-1]
				temp2 := strings.Split(temp, "?")
				pair := temp2[0]

				switch len(msg.Events) {
				case 1:
					m.Lock()
					q := (m.data[pair]).(Quote)
					updateBidAsk(&q, msg.Events[0])
					m.data[pair] = q
					data <- q
					m.Unlock()
				case 2:
					if t := (msg.Events[0]["type"]).(string); t == "trade" {
						m.Lock()
						q := (m.data[pair]).(Quote)
						updatePriceSizeVol(&q, msg.Events[0])
						updateBidAsk(&q, msg.Events[1])
						m.data[pair] = q
						data <- q
						m.Unlock()
					}
				default:
					continue
				}

			}
		}
	}()

	return nil
}

func updatePriceSizeVol(q *Quote, trade map[string]interface{}) {
	q.Price = (trade["price"]).(string)
	q.Size = (trade["amount"]).(string)
	if vol, err := strconv.ParseFloat(q.Volume, 64); err == nil {
		if tradeVol, err := strconv.ParseFloat((trade["amount"]).(string), 64); err == nil {
			q.Volume = strconv.FormatFloat((vol + tradeVol), 'f', -1, 64)
		}
	}
}

func updateBidAsk(q *Quote, change map[string]interface{}) {
	if rem, err := strconv.ParseFloat((change["remaining"]).(string), 64); err == nil {
		if rem > 0 {
			switch side := (change["side"]).(string); side {
			case "bid":
				q.Bid = (change["price"]).(string)
			case "ask":
				q.Ask = (change["price"]).(string)
			}
		}
	}
}
