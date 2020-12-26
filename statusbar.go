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

	Root          *tview.Flex
	AutoCompleter *autoCompleterType
}

func newStatusBar(rowCnt, colCnt int) *statusBar {
	sb := &statusBar{}
	sb.InputText = tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack)
	sb.InputText.SetLabel("Input: ")
	sb.InputText.SetInputCapture(sb.inputTextKeyHandler)

	sb.MessageText = tview.NewTextView().SetScrollable(true).SetTextAlign(tview.AlignLeft)
	sb.CommandText = tview.NewTextView().SetScrollable(true).SetTextAlign(tview.AlignRight)
	sb.RowCountText = tview.NewTextView().SetScrollable(true).SetTextAlign(tview.AlignRight)

	sb.MessageText.SetDynamicColors(true)
	sb.AutoCompleter = newAutoCompleter()

	if sb.AutoCompleter == nil {
		return nil
	}
	// defer ac.destroy()
	sb.InputText.SetAutocompleteFunc(sb.AutoCompleter.getFilteredLines)

	sb.CommandText.SetText("welcome")
	rowCntWidth := len(fmt.Sprintf(" rows[%d,%d]%d", rowCnt, colCnt, rowCnt))

	sb.Root = tview.NewFlex().
		AddItem(sb.InputText, 0, 4, false).
		AddItem(sb.MessageText, 0, 2, false).
		AddItem(sb.CommandText, 10, 1, false).
		AddItem(sb.RowCountText, rowCntWidth, 1, false).
		SetDirection(tview.FlexColumn)
	return sb
}

func (sb *statusBar) getRoot() *tview.Flex {
	return sb.Root
}

func (sb *statusBar) setInputLabel(label string) {
	sb.InputText.SetLabel(label)
}

func (sb *statusBar) setCommandText(str string) {
	sb.CommandText.SetText("(" + str + ")")
}

func (sb *statusBar) setMessageText(msg string) {
	sb.MessageText.SetText(msg)
}

func (sb *statusBar) setRowCountText(count, currentRow, currentCol int) {
	sb.RowCountText.SetText(fmt.Sprintf("rows[%d,%d]%d", currentRow, currentCol, count))
}

func (sb *statusBar) setInputDoneFunc(callback func(cmd, text string) bool) {
	decorator := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cmd := sb.InputText.GetLabel()
			text := sb.InputText.GetText()

			sb.AutoCompleter.addLine(text)
			inputDone := callback(cmd, text)
			if !inputDone {
				return
			}
		}
		sb.InputText.SetText("")
		sb.setInputLabel("Input: ")
	}

	// SetDoneFunc sets a handler which is called when the user is done entering text.
	// The callback function is provided with the key that was pressed, which is one of the following:
	// - KeyEnter: Done entering text.
	// - KeyEscape: Abort text input.
	// - KeyTab: Move to the next field.
	// - KeyBacktab: Move to the previous field.
	sb.InputText.SetDoneFunc(decorator)
}

func (sb *statusBar) inputTextKeyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		sb.AutoCompleter.ReturnAllLines = true
		sb.InputText.Autocomplete()
		sb.AutoCompleter.ReturnAllLines = false
	}
	return event
}
