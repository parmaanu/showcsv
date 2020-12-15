package showcsv

import "errors"

// Display function renders the tui table on the terminal
func Display(tableConfig *TableConfig) error {
	// ac := newAutoCompleter()
	// fmt.Println(ac.addAndGetLines("hello"))
	// fmt.Println(ac.addAndGetLines("how"))
	// fmt.Println(ac.addAndGetLines("are"))
	// return nil

	app := newTuiApp(tableConfig)
	if app == nil {
		return errors.New("Error while creating tuiApp")
	}
	return app.run()
}
