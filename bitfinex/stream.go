package bitfinex

import (
	"errors"
	"math"
	"strconv"

	"github.com/3cb/cq/cq"
	"github.com/gorilla/websocket"
)

// Subscribe is the message structure to subscribe to Bitfinex
// websocket API: https://docs.bitfinex.com/v2/reference#ws-public-ticker
type Subscribe struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
}

// ChannelStore keeps track of channel ID numbers for filtering websocket messages
type ChannelStore struct {
	Trades  map[float64]string
	Tickers map[float64]string
}

func connectWS(m *Market, timerCh chan<- cq.TimerMsg) error {
	var subMsgs []*Subscribe

	m.RLock()
	pairs := m.pairs
	m.RUnlock()
	// create two subscribe messages for each trading pair
	for _, pair := range pairs {
		subTicker := &Subscribe{
			Event:   "subscribe",
			Channel: "ticker",
			Symbol:  pair,
		}
		subMsgs = append(subMsgs, subTicker)
		subTrades := &Subscribe{
			Event:   "subscribe",
			Channel: "trades",
			Symbol:  pair,
		}
		subMsgs = append(subMsgs, subTrades)
	}

	conn, resp, err := websocket.DefaultDialer.Dial("wss://api.bitfinex.com/ws/2", nil)
	if resp.StatusCode != 101 || err != nil {
		return errors.New("unable to connect to bitfinex websocket api")
	}

	go func() {
		store := ChannelStore{
			Trades:  make(map[float64]string),
			Tickers: make(map[float64]string),
		}

		for {
			var data interface{}

			err := conn.ReadJSON(&data)
			if err != nil {
				conn.Close()
				return
			}

			switch msg := data.(type) {
			// this case asserts subsciption confirmation
			case map[string]interface{}:
				registerChannelID(msg, &store)

			// this case asserts snapshots and updates
			case []interface{}:
				id := (msg[0]).(float64)
				switch x := msg[1].(type) {
				// this case asserts ticker snapshot and updates or trades snapshot
				case []interface{}:
					pair, _ := queryStore(&store, id)
					if upd, ok := (x[0]).([]interface{}); ok {
						tradeToQuote(m, upd, pair, timerCh)
					} else {
						tickerToQuote(m, x, pair, timerCh)
					}

				// this case asserts trades updates
				case string:
					if msg[1] == "tu" {
						upd := (msg[2]).([]interface{})
						pair, _ := queryStore(&store, id)
						tradeToQuote(m, upd, pair, timerCh)
					} // handle "hb" case here ***
				}
			}
		}
	}()

	// send subscribe messages
	for _, msg := range subMsgs {
		conn.WriteJSON(msg)
	}

	return nil
}

func registerChannelID(msg map[string]interface{}, store *ChannelStore) {
	event := (msg["event"]).(string)
	if event == "subscribed" {
		chanID := (msg["chanId"]).(float64)
		symbol := (msg["symbol"]).(string)
		switch channel := (msg["channel"]).(string); channel {
		case "ticker":
			store.Tickers[chanID] = symbol
		case "trades":
			store.Trades[chanID] = symbol
		}
	}
}

// queryStore uses channel id to get pair string
// -> first return value is pair string
// -> second return value is channel type: "trade" or "ticker"
func queryStore(store *ChannelStore, id float64) (string, string) {
	if pair, ok := store.Trades[id]; ok {
		return pair, "trade"
	}
	pair := store.Tickers[id]
	return pair, "ticker"
}

func tickerToQuote(m *Market, upd []interface{}, pair string, timerCh chan<- cq.TimerMsg) {
	m.Lock()
	q := m.data[pair]
	q.Bid = strconv.FormatFloat((upd[0]).(float64), 'f', -1, 64)
	q.Ask = strconv.FormatFloat((upd[2]).(float64), 'f', -1, 64)
	q.Change = strconv.FormatFloat((upd[4]).(float64), 'f', -1, 64)
	q.ChangePerc = strconv.FormatFloat((upd[5]).(float64), 'f', -1, 64)
	q.Volume = strconv.FormatFloat((upd[7]).(float64), 'f', -1, 64)
	q.High = strconv.FormatFloat((upd[8]).(float64), 'f', -1, 64)
	q.Low = strconv.FormatFloat((upd[9]).(float64), 'f', -1, 64)
	m.data[pair] = q
	m.Unlock()
	timerCh <- cq.TimerMsg{
		IsTrade: false,
		Quote:   q,
	}
}

func tradeToQuote(m *Market, upd []interface{}, pair string, timerCh chan<- cq.TimerMsg) {
	m.Lock()
	q := m.data[pair]
	q.Size = strconv.FormatFloat(math.Abs((upd[2]).(float64)), 'f', -1, 64)
	q.Price = strconv.FormatFloat((upd[3]).(float64), 'f', -1, 64)
	m.data[pair] = q
	m.Unlock()
	timerCh <- cq.TimerMsg{
		IsTrade: true,
		Quote:   q,
	}
}
