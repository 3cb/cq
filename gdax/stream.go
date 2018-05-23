package gdax

import (
	"errors"

	"github.com/3cb/cq/display"
	"github.com/gorilla/websocket"
)

// Subscribe is the structure for the subscription message sent to GDAX websocket API
// https://docs.gdax.com/#subscribe
type Subscribe struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

func connectWS(m *Market, upd chan display.Setter) error {
	wsSub := &Subscribe{
		Type:       "subscribe",
		ProductIds: m.pairs,
		Channels:   []string{"matches", "ticker"},
	}

	conn, resp, err := websocket.DefaultDialer.Dial("wss://ws-feed.gdax.com", nil)
	if resp.StatusCode != 101 || err != nil {
		return errors.New("unable to connect to gdax websocket api")
	}
	conn.WriteJSON(wsSub)

	go func() {
		msg := Quote{}
		for {
			err := conn.ReadJSON(&msg)
			if err != nil {
				conn.Close()
				// handle error here
				return
			}

			if msg.Type == "match" {
				m.Lock()
				quote := m.data[msg.ID]
				quote.Price = msg.Price
				quote.Size = msg.Size
				m.data[msg.ID] = quote
				upd <- quote
				m.Unlock()
			} else if msg.Type == "ticker" {
				m.Lock()
				quote := m.data[msg.ID]
				quote.Bid = msg.Bid
				quote.Ask = msg.Ask
				quote.High = msg.High
				quote.Low = msg.Low
				quote.Open = msg.Open
				quote.Volume = msg.Volume
				m.data[msg.ID] = quote
				upd <- quote
				m.Unlock()
			}
			msg = Quote{}
		}
	}()

	return nil
}
