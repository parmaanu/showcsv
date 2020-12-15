package showcsv

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tuiApp struct {
	MainApp   *tview.Application
	TableView *tview.Table
	StatusBar *statusBar

	totalRowCnt int
	totalColCnt int
	currentRow  int
	currentCol  int

	defaultStyle   tcell.Style
	highlightStyle tcell.Style

	inputMode          bool
	previousSearchText string
}

func newTuiApp(tableConfig *TableConfig) *tuiApp {
	app := &tuiApp{}

	app.MainApp = tview.NewApplication()
	app.MainApp.SetInputCapture(app.applicationInput)

	app.totalRowCnt = len(tableConfig.Data)
	if len(tableConfig.Data) > 0 {
		app.totalColCnt = len(tableConfig.Data[0])
	}

	app.TableView = newTableView(tableConfig)
	app.StatusBar = newStatusBar(app.totalRowCnt, app.totalColCnt)
	if app.StatusBar == nil {
		return nil
	}
	app.StatusBar.setRowCountText(app.totalRowCnt, app.currentRow, app.currentCol+1)
	// app.StatusBar.InputText.SetDoneFunc(app.userInputDoneFunc)
	app.StatusBar.setInputDoneFunc(app.userInputDoneFunc)

	app.defaultStyle = tcell.Style{}.Attributes(tcell.AttrNone).Foreground(tcell.ColorWhite)
	app.highlightStyle = tcell.Style{}.Attributes(tcell.AttrBold).Foreground(tcell.ColorWhite)

	appInitialized := false
	app.MainApp.SetAfterDrawFunc(func(screen tcell.Screen) {
		if !appInitialized {
			appInitialized = true
			app.hightlightRow(1, app.highlightStyle)
			app.hightlightCol(0, app.highlightStyle)
		}
	})

	app.TableView.SetSelectionChangedFunc(app.tableSelectionChangedCallback)
	return app
}

func (app *tuiApp) run() error {
	if app.MainApp == nil || app.TableView == nil || app.StatusBar.InputText == nil {
		return errors.New("tuiApp not initialized properly")
	}
	flex := tview.NewFlex().
		AddItem(app.TableView, 0, 1, true).
		AddItem(app.StatusBar.getRoot(), 1, 1, false).
		SetDirection(tview.FlexRow)
	return app.MainApp.SetRoot(flex, true).EnableMouse(true).Run()
}

func (app *tuiApp) userInputDoneFunc(cmd, text string) {
	switch cmd {
	case "/", " ":
		app.searchColForward(text)
	}
	app.inputMode = false
	app.MainApp.SetFocus(app.TableView)
}

func (app *tuiApp) applicationInput(event *tcell.EventKey) *tcell.EventKey {
	// TODO, don't set command when in command mode
	// TODO, don't quit application on `q` when in command mode
	app.StatusBar.setCommandText(event.Name())
	app.StatusBar.setMessageText("")
	switch event.Key() {
	case tcell.KeyCtrlQ:
		app.MainApp.Stop()
	case tcell.KeyCtrlS:
		// TODO, save
	case tcell.KeyCtrlH:
		// TODO, show help
		// case tcell.
	case tcell.KeyCtrlC:
		// TODO, cancel all pending/remaining tasks or threads
		app.StatusBar.setInputLabel("Input: ")
		app.MainApp.SetFocus(app.TableView)
		app.StatusBar.InputText.SetText("")
		app.inputMode = false
		return nil
	case tcell.KeyEsc:
		app.StatusBar.setInputLabel("Input: ")
		app.MainApp.SetFocus(app.TableView)
		return nil
	case tcell.KeyRune:
		switch ch := event.Rune(); ch {
		case 'q', 'Q':
			if !app.inputMode {
				app.MainApp.Stop()
			}
		case 'n', 'N':
			// TODO, go to next search item
			if len(app.previousSearchText) > 0 {
				if ch == 'n' {
					app.searchColForward(app.previousSearchText)
				} else if ch == 'N' {
					app.searchColBackward(app.previousSearchText)
				}
			} else {
				app.StatusBar.setMessageText(fmt.Sprintf("no previous search pattern!"))
			}
		case ' ':
			ch = '/'
			fallthrough
		case '/', '?':
			if !app.inputMode {
				app.inputMode = true
				app.StatusBar.setInputLabel(string(ch))
				app.MainApp.SetFocus(app.StatusBar.InputText)
				// return nil
			}
		case ':':
			// TODO, split mode
			if !app.inputMode {
				app.inputMode = true
				app.StatusBar.setInputLabel("split-regex:")
				app.MainApp.SetFocus(app.StatusBar.InputText)
				return nil
			}
		case '|':
			// TODO, select mode
			if !app.inputMode {
				app.inputMode = true
				app.StatusBar.setInputLabel("select-regex:")
				app.MainApp.SetFocus(app.StatusBar.InputText)
				return nil
			}
		case ',':
			// TODO, select all matching pattern in current column
		case '_':
			// TODO, optimize column width
		case '-':
			// TODO, hide column
		case 'K':
			// TODO, move row up
		case 'J':
			// TODO, move row down
		case 'H':
			// TODO, move column left
		case 'L':
			// TODO, move column right
		}
	}
	return event
}

func (app *tuiApp) searchColForward(searchText string) {
	if len(searchText) == 0 {
		return
	}
	app.previousSearchText = searchText
	searchFunc := func(s string) bool { return strings.Contains(s, strings.ToLower(searchText)) }
	reg, err := regexp.Compile(searchText)
	if err == nil {
		searchFunc = func(s string) bool { return reg.MatchString(s) }
	}

	startRow, col := app.TableView.GetSelection()
	for row := startRow; row < app.totalRowCnt; row++ {
		cellText := app.TableView.GetCell(row, col).Text
		if searchFunc(cellText) {
			app.TableView.Select(row, col)
			app.StatusBar.setMessageText(fmt.Sprintf("found /%s/", searchText))
			return
		}
	}
	for row := 1; row < startRow; row++ {
		cellText := app.TableView.GetCell(row, col).Text
		if searchFunc(cellText) {
			app.TableView.Select(row, col)
			app.StatusBar.setMessageText(fmt.Sprintf("found from top /%s/", searchText))
			return
		}
	}
	app.StatusBar.setMessageText(fmt.Sprintf("not found [red]/%s/", searchText))
}

func (app *tuiApp) searchColBackward(searchText string) {
	if len(searchText) == 0 {
		return
	}
	app.previousSearchText = searchText
	searchFunc := func(s string) bool { return strings.Contains(s, strings.ToLower(searchText)) }
	reg, err := regexp.Compile(searchText)
	if err == nil {
		searchFunc = func(s string) bool { return reg.MatchString(s) }
	}

	startRow, col := app.TableView.GetSelection()
	for row := startRow; row > 0; row-- {
		cellText := app.TableView.GetCell(row, col).Text
		if searchFunc(cellText) {
			app.TableView.Select(row, col)
			app.StatusBar.setMessageText(fmt.Sprintf("found /%s/", searchText))
			return
		}
	}
	for row := app.totalRowCnt; row > startRow; row-- {
		cellText := app.TableView.GetCell(row, col).Text
		if searchFunc(cellText) {
			app.TableView.Select(row, col)
			app.StatusBar.setMessageText(fmt.Sprintf("found from bottom /%s/", searchText))
			return
		}
	}
	app.StatusBar.setMessageText(fmt.Sprintf("not found [red]/%s/", searchText))
}

func (app *tuiApp) tableSelectionChangedCallback(row, col int) {
	// Note, when first row is not selectable then row is -1, similarly for col
	if row < 0 || col < 0 {
		return
	}

	// Note, set everything on visible screen to defaultStyle.
	// If you set previous row, col to defaultStyle then there is a problem in hightlighting
	// while jumping out of visible screen.
	app.resetScreenStyle(app.defaultStyle)
	app.hightlightRow(row, app.highlightStyle)
	app.hightlightCol(col, app.highlightStyle)

	app.currentRow = row
	app.currentCol = col
	app.StatusBar.setRowCountText(app.totalRowCnt, app.currentRow, app.currentCol+1)
}

func (app *tuiApp) hightlightRow(row int, style tcell.Style) {
	_, startCol := app.TableView.GetOffset()
	_, _, tableWidth, _ := app.TableView.Box.GetRect()
	endCol := startCol + tableWidth
	if endCol > app.totalColCnt {
		endCol = app.totalColCnt
	}
	for col := startCol; col < endCol; col++ {
		app.TableView.GetCell(row, col).SetStyle(style)
	}
}

func (app *tuiApp) hightlightCol(col int, style tcell.Style) {
	startRow, _ := app.TableView.GetOffset()
	_, _, _, tableHeight := app.TableView.Box.GetRect()
	endRow := startRow + tableHeight
	if endRow > app.totalRowCnt+1 {
		endRow = app.totalRowCnt + 1
	}

	// Note, row starts with startRow + 1 as first row is header
	for row := startRow + 1; row < endRow; row++ {
		app.TableView.GetCell(row, col).SetStyle(style)
	}
}

func (app *tuiApp) resetScreenStyle(style tcell.Style) {
	_, _, tableWidth, tableHeight := app.TableView.Box.GetRect()
	startRow, startCol := app.TableView.GetOffset()
	endRow := startRow + tableHeight
	endCol := startCol + tableWidth

	if endRow > app.totalRowCnt+1 {
		endRow = app.totalRowCnt + 1
	}
	if endCol > app.totalColCnt {
		endCol = app.totalColCnt
	}

	// Note, row starts with startRow + 1 as first row is header
	for row := startRow + 1; row < endRow; row++ {
		for col := startCol; col < endCol; col++ {
			app.TableView.GetCell(row, col).SetStyle(style)
		}
	}
}
