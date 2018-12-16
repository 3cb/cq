package cq

import "github.com/gdamore/tcell"

// Config contains ui style settings
type Config struct {
	ID    string
	Theme struct {
		// main colors
		ShowTableBorders bool
		BackgroundColor  tcell.Color
		BorderColor      tcell.Color
		HeaderTextColor  tcell.Color
		TextColor        tcell.Color
		TextColorUp      tcell.Color
		TextColorDown    tcell.Color

		// menu colors
		TextColorMenuMain           tcell.Color
		TextColorMenuSecondary      tcell.Color
		TextColorMenuShortcut       tcell.Color
		TextColorMenuSelected       tcell.Color
		BackgroundColorMenuSelected tcell.Color
	}
	CellFlash tcell.AttrMask
}

// SetConfig sets the ui style configuration
func SetConfig(t *string, f *bool) Config {
	c := Config{}
	switch *t {
	case "light":
		c.ID = "light"

		// set main colors
		c.Theme.ShowTableBorders = true
		c.Theme.BackgroundColor = tcell.Color255
		c.Theme.BorderColor = tcell.ColorLightSlateGrey
		c.Theme.HeaderTextColor = tcell.ColorBlue
		c.Theme.TextColor = tcell.Color238
		c.Theme.TextColorUp = tcell.ColorGreen
		c.Theme.TextColorDown = tcell.ColorRed

		// set menu colors
		c.Theme.TextColorMenuMain = tcell.Color238
		c.Theme.TextColorMenuSecondary = tcell.ColorGreen
		c.Theme.TextColorMenuShortcut = tcell.ColorBlue
		c.Theme.TextColorMenuSelected = tcell.Color255
		c.Theme.BackgroundColorMenuSelected = tcell.Color237
	// case: "dark":
	default:
		c.ID = "dark"

		// set main colors
		c.Theme.ShowTableBorders = true
		c.Theme.BackgroundColor = tcell.ColorBlack
		c.Theme.BorderColor = tcell.ColorLightSlateGray
		c.Theme.HeaderTextColor = tcell.ColorYellow
		c.Theme.TextColor = tcell.ColorWhite
		c.Theme.TextColorUp = tcell.ColorGreen
		c.Theme.TextColorDown = tcell.ColorRed

		// set menu colors
		c.Theme.TextColorMenuMain = tcell.ColorWhite
		c.Theme.TextColorMenuSecondary = tcell.ColorGreen
		c.Theme.TextColorMenuShortcut = tcell.ColorYellow
		c.Theme.TextColorMenuSelected = tcell.ColorBlack
		c.Theme.BackgroundColorMenuSelected = tcell.ColorWhite
	}

	switch *f {
	case false:
		c.CellFlash = tcell.AttrBold
	default:
		c.CellFlash = tcell.AttrReverse
	}

	return c
}
