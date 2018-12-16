package cq

import (
	"sync"
	"time"
)

const queueSize = 200
const timerDuration = (300 * time.Millisecond)

// TimerGroup contains matrix of all the flash timers
type TimerGroup struct {
	sync.RWMutex

	// exchange | pair | channel
	list map[string]map[string]chan TimerMsg
}

// StartTimerGroup creates double map with timers
func StartTimerGroup(exchanges map[string]Exchange) (<-chan UpdateMsg, chan<- TimerMsg) {
	tg := &TimerGroup{
		list: make(map[string]map[string]chan TimerMsg),
	}
	updateCh := make(chan UpdateMsg, queueSize)
	timerCh := make(chan TimerMsg, queueSize)

	for k, e := range exchanges {
		tg.list[k] = make(map[string]chan TimerMsg)

		for _, id := range e.GetIDs() {
			tg.list[k][id] = make(chan TimerMsg, queueSize)
		}
	}

	go launchTimers(tg, updateCh, timerCh)

	return updateCh, timerCh
}

// start goroutine per exchange per trading pair to track flash times
func launchTimers(tg *TimerGroup, updateCh chan<- UpdateMsg, routerCh <-chan TimerMsg) {
	// start individual timer loops
	tg.Lock()
	for _, pairMap := range tg.list {
		for _, ch := range pairMap {
			go func(ch <-chan TimerMsg) {
				var lastTime time.Time
				var lastQuote Quote

				timer := time.NewTimer(timerDuration)

				// ignore first value from timer
				<-timer.C

				for {
					select {
					case t := <-timer.C:
						if t.After(lastTime) {
							updateCh <- UpdateMsg{
								IsTrade: true,
								Flash:   false,
								Quote:   lastQuote,
							}
						}
					case msg := <-ch:
						switch msg.IsTrade {
						case true:
							timer.Stop()
							timer.Reset(timerDuration)
							lastTime = time.Now()
							lastQuote = msg.Quote
							updateCh <- UpdateMsg{
								Quote:   msg.Quote,
								IsTrade: true,
								Flash:   true,
							}
						case false:
							lastQuote = msg.Quote
							updateCh <- UpdateMsg{
								Quote:   msg.Quote,
								IsTrade: false,
								Flash:   false,
							}
						}
					}
				}
			}(ch)
		}
	}
	tg.Unlock()

	// router loop directs messages to corresponding timer loop
	for {
		msg := <-routerCh
		tg.RLock()
		ch := tg.list[msg.Quote.MarketID][msg.Quote.ID]
		tg.RUnlock()
		ch <- msg
	}
}
