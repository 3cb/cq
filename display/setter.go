package display

import "github.com/rivo/tview"

// Setter is an interface that market quotes implement to initialize
// and update data in gui table
type Setter interface {
	SetRow(*tview.Table)
	UpdRow(*tview.Table)
	ClrBold(*tview.Table)
}
