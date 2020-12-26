package showcsv

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO, tableview should allow inserting a row and a column
// TODO, editable cells
// TODO, Add a functionality to add numeric column - InsertColumn(0) and tablecell.SetSelectable(false)

// TableConfig for tuitable
type TableConfig struct {
	// TODO, make header editable to change the column name
	Name   string
	Header []string
	Data   [][]string
}

func newTableView(tableConfig *TableConfig) *tview.Table {
	table := tview.NewTable().
		SetFixed(1, 0)

	table.Select(0, 0).
		SetSelectable(true, true)

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

	return table
}
