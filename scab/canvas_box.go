package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CanvasBox struct {
	*tview.Box
	canvas  *Canvas
	cursorX *int
	cursorY *int
	showGrid *bool
}

func (c *CanvasBox) SetCanvas(canvas *Canvas) {
	c.canvas = canvas
}

func NewCanvasBox(canvas *Canvas, cursorX, cursorY *int, showGrid *bool) *CanvasBox {
	return &CanvasBox{
		Box:      tview.NewBox(),
		canvas:   canvas,
		cursorX:  cursorX,
		cursorY:  cursorY,
		showGrid: showGrid,
	}
}

func (c *CanvasBox) Draw(screen tcell.Screen) {
	c.Box.DrawForSubclass(screen, c)

	x, y, w, h := c.Box.GetInnerRect()
	if w <= 0 || h <= 0 {
		return
	}

	canvasBG := tcell.ColorDefault
	if c.canvas.BGColor != "" {
		canvasBG = tcell.GetColor(c.canvas.BGColor)
	}
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			screen.SetContent(x+dx, y+dy, ' ', nil, tcell.StyleDefault.Background(canvasBG))
		}
	}

	maxY := c.canvas.Height
	if maxY > h {
		maxY = h
	}
	maxX := c.canvas.Width
	if maxX > w {
		maxX = w
	}

	for canvasY := 0; canvasY < maxY; canvasY++ {
		for canvasX := 0; canvasX < maxX; canvasX++ {
			cell := c.canvas.GetCell(canvasX, canvasY)
			if cell == nil {
				continue
			}

			ch := ' '
			if cell.Char != "" {
				runes := []rune(cell.Char)
				if len(runes) > 0 {
					ch = runes[0]
				}
			}

			fg := tcell.GetColor(cell.FG)
			bg := canvasBG
			if cell.BG != "" {
				bg = tcell.GetColor(cell.BG)
			}

			if *c.showGrid && ch == ' ' && cell.BG == "" {
				bg = tcell.GetColor("#3d3d5c")
			}

			if canvasX == *c.cursorX && canvasY == *c.cursorY {
				fg = tcell.GetColor("#000000")
				bg = tcell.GetColor("#ffffff")
			}

			style := tcell.StyleDefault.Foreground(fg).Background(bg)
			screen.SetContent(x+canvasX, y+canvasY, ch, nil, style)
		}
	}
}
