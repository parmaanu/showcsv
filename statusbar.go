package showcsv

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TODO, show `filename | ctrl-H opens help`` in status bar initially on the start
// TODO, then show filename | message									button pressed button-function		rows
type statusBar struct {
	InputText    *tview.InputField // TODO, Can InputText and MessageText be merged - only one of them would be used at a time?
	MessageText  *tview.TextView   // shows error messages, alerts etc to user
	CommandText  *tview.TextView   // shows key pressed and commands
	RowCountText *tview.TextView   // shows the number of rows in the sheet

	Root *tview.Flex
}

func newStatusBar(rowCnt, colCnt int) *statusBar {
	sb := &statusBar{}
	sb.InputText = tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack)
	sb.InputText.SetLabel("Input: ")

	sb.MessageText = tview.NewTextView().SetScrollable(true).SetTextAlign(tview.AlignLeft)
	sb.CommandText = tview.NewTextView().SetScrollable(true).SetTextAlign(tview.AlignRight)
	sb.RowCountText = tview.NewTextView().SetScrollable(true).SetTextAlign(tview.AlignRight)

	sb.CommandText.SetText("welcome")
	rowCntWidth := len(fmt.Sprintf(" rows[%d,%d]%d", rowCnt, colCnt, rowCnt))

	// TODO, tweak the sizes (fixed of flexible) of right 3 text boxes
	sb.Root = tview.NewFlex().
		AddItem(sb.InputText, 0, 4, false).
		AddItem(sb.MessageText, 0, 2, false).
		AddItem(sb.CommandText, 10, 1, false).
		AddItem(sb.RowCountText, rowCntWidth, 1, false).
		SetDirection(tview.FlexColumn)

	// TODO, set autocomplete in InputText

	// statusBar.SetBackgroundColor(tcell.ColorGray)
	// SetRect()
	// TODO, do SetBorder when highlighted. Do I need to use SetRect(x, y, width, height)

	// app.StatusBar.SetDoneFunc(func(key tcell.Key) {
	//     // statusBar..
	//     text := app.StatusBar.GetText()
	//     text = strings.TrimLeft(text, "Input: ")
	//     app.StatusBar.SetText("Input: " + text)
	// })
	return sb
}

func (sb *statusBar) getRoot() *tview.Flex {
	return sb.Root
}

func (sb *statusBar) setCommandText(str string) {
	sb.CommandText.SetText("(" + str + ")")
}

func (sb *statusBar) setRowCountText(count, currentRow, currentCol int) {
	sb.RowCountText.SetText(fmt.Sprintf("rows[%d,%d]%d", currentRow, currentCol, count))
}

func (sb *statusBar) setMessageText(msg string) {
	sb.MessageText.SetText(msg)
}
