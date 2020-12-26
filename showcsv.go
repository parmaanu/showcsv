package showcsv

import (
	"errors"
	"os"
)

// Display function renders the tui table on the terminal
// Please note that you have to set SHOWCSV_RENDER_TUI environment variable to render TUI.
// `export SHOWCSV_RENDER_TUI=true`
// You can save above line in your ~/.bashrc to make it persistent
func Display(tableConfig *TableConfig) error {
	_, renderTui := os.LookupEnv("SHOWCSV_RENDER_TUI")
	if !renderTui {
		return nil
	}

	app := newTuiApp(tableConfig)
	if app == nil {
		return errors.New("Error while creating tuiApp")
	}
	return app.run()
}
