package showcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
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

// DisplayFile loads the given csvfile and displays it on the terminal
func DisplayFile(fname string, firstRowIsHeader bool) error {
	csvFile, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open csv file", err)
	}

	r := csv.NewReader(csvFile)
	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("error while reading records from csv", err)
	}

	if len(records) == 0 {
		return errors.New("empty csv file")
	}

	header := make([]string, len(records[0]))
	if firstRowIsHeader {
		header = records[0]
		records = records[1:]
	} else {
		for i := 0; i < len(records[0]); i++ {
			header[i] = fmt.Sprintf("Col%d", i+1)
		}
	}

	tableConfig := &TableConfig{fname, header, records}
	return Display(tableConfig)
}
