package showcsv

import (
	"errors"
)

// Display function renders the tui table on the terminal
func Display(tableConfig *TableConfig) error {
	app := newTuiApp(tableConfig)
	if app == nil {
		return errors.New("Error while creating tuiApp")
	}
	return app.render()
}
