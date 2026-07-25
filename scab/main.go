package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeErase
	ModeBounds
)

type App struct {
	app        *tview.Application
	canvas     *Canvas
	pages      *tview.Pages
	canvasBox  *CanvasBox
	statusBar  *tview.TextView
	canvasGrid *tview.Grid

	mode       Mode
	cursorX    int
	cursorY    int
	scale      int
	showGrid   bool
	lastChar   rune
	canvasName string
	dirty      bool

	boundsOriginX int
	boundsOriginY int
}

func NewApp() *App {
	a := &App{
		app:        tview.NewApplication(),
		canvas:     NewCanvas(40, 20, "untitled"),
		pages:      tview.NewPages(),
		statusBar:  tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft),
		mode:       ModeNormal,
		scale:      1,
		showGrid:   false,
		canvasName: "untitled",
	}

	a.setupUI()
	return a
}

func (a *App) setupUI() {
	a.updateStatusBar()

	a.canvasBox = NewCanvasBox(a.canvas, &a.cursorX, &a.cursorY, &a.showGrid)
	a.canvasBox.SetBorder(true)
	a.canvasBox.SetBorderColor(tcell.GetColor("#6c6c8a"))
	a.canvasBox.SetTitle(" Canvas ")

	a.canvasGrid = tview.NewGrid().
		SetRows(0, a.canvas.Height+2, 0).
		SetColumns(0, a.canvas.Width+2, 0).
		AddItem(a.canvasBox, 1, 1, 1, 1, 0, 0, true)
	a.canvasGrid.SetBackgroundColor(tcell.GetColor("#14141e"))

	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.canvasGrid, 0, 5, true).
		AddItem(a.statusBar, 1, 1, false)

	a.pages.AddPage("main", rootLayout, true, true)

	a.app.SetRoot(a.pages, true)
	a.app.SetInputCapture(a.handleInput)
}

func (a *App) modeName() string {
	switch a.mode {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeErase:
		return "ERASE"
	case ModeBounds:
		return "BOUNDS"
	}
	return "?"
}

func (a *App) updateStatusBar() {
	modeColor := "white"
	switch a.mode {
	case ModeInsert:
		modeColor = "green"
	case ModeErase:
		modeColor = "red"
	case ModeBounds:
		modeColor = "yellow"
	}

	dirtyMark := ""
	if a.dirty {
		dirtyMark = " *"
	}

	status := fmt.Sprintf(" [%s]%s[white] %s | (%d,%d) | %dx%d | g:%v",
		modeColor, a.modeName(), dirtyMark,
		a.cursorX, a.cursorY,
		a.canvas.Width, a.canvas.Height,
		a.showGrid)
	a.statusBar.SetText(status)
}

func (a *App) renderCanvas() {
	a.app.ForceDraw()
}

func (a *App) applyAction() {
	switch a.mode {
	case ModeInsert:
		char := string(a.lastChar)
		if char == "" {
			char = " "
		}
		a.canvas.SetCell(a.cursorX, a.cursorY, char)
		a.dirty = true
	case ModeErase:
		a.canvas.SetCell(a.cursorX, a.cursorY, " ")
		a.dirty = true
	}
}

func (a *App) moveCursor(dx, dy int) {
	newX := a.cursorX + dx
	newY := a.cursorY + dy

	if newX >= 0 && newX < a.canvas.Width && newY >= 0 && newY < a.canvas.Height {
		a.cursorX = newX
		a.cursorY = newY
		a.renderCanvas()
		a.updateStatusBar()
	}
}

func (a *App) resizeCanvas(deltaW, deltaH int) {
	newW := a.canvas.Width + deltaW
	newH := a.canvas.Height + deltaH
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}
	if newW > 200 {
		newW = 200
	}
	if newH > 100 {
		newH = 100
	}
	a.canvas.Resize(newW, newH)
	if a.cursorX >= newW {
		a.cursorX = newW - 1
	}
	if a.cursorY >= newH {
		a.cursorY = newH - 1
	}
	a.canvasGrid.Clear().SetRows(0, newH+2, 0).SetColumns(0, newW+2, 0).
		AddItem(a.canvasBox, 1, 1, 1, 1, 0, 0, true)
	a.dirty = true
	a.renderCanvas()
	a.updateStatusBar()
}

func (a *App) handleInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEscape {
		if a.pages.HasPage("dialog") {
			a.pages.RemovePage("dialog")
		}
		a.mode = ModeNormal
		a.renderCanvas()
		a.updateStatusBar()
		return nil
	}

	if a.pages.HasPage("dialog") {
		return event
	}

	if a.mode == ModeBounds {
		return a.handleBoundsInput(event)
	}

	if a.mode == ModeNormal {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'i':
				a.mode = ModeInsert
				a.updateStatusBar()
				return nil
			case 'e':
				a.mode = ModeErase
				a.updateStatusBar()
				return nil
			case 'b':
				a.mode = ModeBounds
				a.boundsOriginX = a.canvas.Width
				a.boundsOriginY = a.canvas.Height
				a.updateStatusBar()
				return nil
			case 'l':
				a.showLoadDialog()
				return nil
			case 's':
				a.saveCanvas()
				return nil
			case 'S':
				a.showSaveAsDialog()
				return nil
			case 'g':
				a.showGrid = !a.showGrid
				a.renderCanvas()
				a.updateStatusBar()
				return nil
			}
		}
	}

	if a.mode == ModeInsert || a.mode == ModeErase {
		switch event.Key() {
		case tcell.KeyUp:
			a.moveCursor(0, -1)
			return nil
		case tcell.KeyDown:
			a.moveCursor(0, 1)
			return nil
		case tcell.KeyLeft:
			a.moveCursor(-1, 0)
			return nil
		case tcell.KeyRight:
			a.moveCursor(1, 0)
			return nil
		case tcell.KeyRune:
			a.lastChar = event.Rune()
			a.applyAction()
			a.renderCanvas()
			return nil
		}
	}

	if a.mode == ModeNormal {
		switch event.Key() {
		case tcell.KeyUp:
			a.moveCursor(0, -1)
			return nil
		case tcell.KeyDown:
			a.moveCursor(0, 1)
			return nil
		case tcell.KeyLeft:
			a.moveCursor(-1, 0)
			return nil
		case tcell.KeyRight:
			a.moveCursor(1, 0)
			return nil
		}
	}

	return event
}

func (a *App) handleBoundsInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyUp:
		a.resizeCanvas(0, -1)
		return nil
	case tcell.KeyDown:
		a.resizeCanvas(0, 1)
		return nil
	case tcell.KeyLeft:
		a.resizeCanvas(-1, 0)
		return nil
	case tcell.KeyRight:
		a.resizeCanvas(1, 0)
		return nil
	}
	return event
}

func (a *App) saveCanvas() {
	a.canvas.Name = a.canvasName
	err := SaveCanvas(a.canvas)
	if err != nil {
		a.statusBar.SetText(fmt.Sprintf(" [red]Save Error: %v[white] ", err))
	} else {
		a.dirty = false
		a.statusBar.SetText(fmt.Sprintf(" [green]Saved: %s[white] ", a.canvasName))
	}
}

func (a *App) showLoadDialog() {
	components := ListComponents()
	if len(components) == 0 {
		a.statusBar.SetText(" [yellow]No components found[white] ")
		a.app.ForceDraw()
		return
	}

	list := tview.NewList()
	for _, name := range components {
		n := name
		list.AddItem(name, "", 0, func() {
			canvas, err := LoadCanvas(n)
			if err == nil {
				a.canvas = canvas
				a.canvasName = n
				a.cursorX = 0
				a.cursorY = 0
				a.dirty = false
				a.canvasBox.SetCanvas(canvas)
				a.canvasGrid.Clear().SetRows(0, canvas.Height+2, 0).SetColumns(0, canvas.Width+2, 0).
					AddItem(a.canvasBox, 1, 1, 1, 1, 0, 0, true)
				a.renderCanvas()
				a.updateStatusBar()
			}
			a.pages.RemovePage("dialog")
		})
	}

	list.SetDoneFunc(func() {
		a.pages.RemovePage("dialog")
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			a.pages.RemovePage("dialog")
			return nil
		}
		return event
	})

	list.SetBackgroundColor(tcell.GetColor("#2d2d44"))
	list.SetBorderColor(tcell.GetColor("#6c6c8a"))
	list.SetBorder(true)
	list.SetTitle(" Load ")

	a.pages.AddPage("dialog", tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(list, 40, 1, true).
		AddItem(nil, 0, 1, false), true, true)
}

func (a *App) showSaveAsDialog() {
	input := tview.NewInputField()
	input.SetLabel("Name: ")
	input.SetText(a.canvasName)
	input.SetFieldWidth(30)
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			name := input.GetText()
			if name != "" {
				a.canvasName = name
				a.canvas.Name = name
				a.saveCanvas()
			}
			a.pages.RemovePage("dialog")
		}
		if key == tcell.KeyEscape {
			a.pages.RemovePage("dialog")
		}
	})

	input.SetBackgroundColor(tcell.GetColor("#2d2d44"))
	input.SetBorderColor(tcell.GetColor("#6c6c8a"))
	input.SetBorder(true)
	input.SetTitle(" Save As ")

	a.pages.AddPage("dialog", tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(input, 40, 1, true).
		AddItem(nil, 0, 1, false), true, true)
}

func main() {
	app := NewApp()

	if err := app.app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running SCAB: %v\n", err)
		os.Exit(1)
	}
}
