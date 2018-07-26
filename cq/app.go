package cq

import (
	"sync"

	"github.com/rivo/tview"
)

// App wraps a tview.Application type with a mutex
type App struct {
	sync.Mutex
	a *tview.Application
}

// Draw calls the Draw() method on internal tview.Application
func (app *App) Draw() {
	app.Lock()
	app.a.Draw()
	app.Unlock()
}
