package showcsv

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO, tableview should allow inserting a row and a column
// TODO, editable cells

// TableConfig for tuitable
type TableConfig struct {
	// TODO, fix the header - cursor should not go to first header? or it can be used to edit column names
	Name   string
	Header []string
	Data   [][]string
}

func newTableView(tableConfig *TableConfig) *tview.Table {
	table := tview.NewTable().
		SetFixed(1, 0)

	table.Select(0, 0).
		SetSelectable(true, true)
		// TODO, fix any one column separator or no separator
		// SetSeparator('|').
		// SetSeparator(tview.Borders.Vertical)

	// TODO, Add a functionality to add numeric column - InsertColumn(0) and tablecell.SetSelectable(false)
	// add header
	tableRowidx := 0
	if len(tableConfig.Header) > 0 {
		headerStyle := tcell.Style{}.Background(tcell.ColorLightGrey).Foreground(tcell.ColorBlack).Bold(true)
		for colidx, headerCol := range tableConfig.Header {
			cell := tview.NewTableCell(headerCol).
				SetAlign(tview.AlignLeft).
				SetStyle(headerStyle).
				SetSelectable(false)
			table.SetCell(tableRowidx, colidx, cell)
		}
		// increment tablew row idx as we have added a header
		tableRowidx++
	}

	// set data in the table
	for rowidx, row := range tableConfig.Data {
		row = tableConfig.Data[rowidx]
		for colidx, ele := range row {
			table.SetCell(tableRowidx, colidx,
				tview.NewTableCell(ele).
					SetAlign(tview.AlignLeft))
		}
		tableRowidx++
	}

	// TODO, highlight background white of column whenever that column is selected
	// TODO, make selected column bold
	// TODO, type / should highlight statusBar
	// TODO, show `filename | ctrl-H opens help`` in status bar initially on the start
	// TODO, then show filename | message									button pressed button-function		rows

	return table
}
