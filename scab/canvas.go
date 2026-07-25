package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Cell struct {
	Char      string `json:"char"`
	FG        string `json:"fg"`
	BG        string `json:"bg"`
	Bold      bool   `json:"bold"`
	Italic    bool   `json:"italic"`
	Underline bool   `json:"underline"`
}

type Canvas struct {
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Cells     [][]*Cell `json:"cells"`
	BGColor   string    `json:"bg_color"`
	Name      string    `json:"name"`
	LatticeOn bool      `json:"lattice_on"`
}

func NewCanvas(width, height int, name string) *Canvas {
	cells := make([][]*Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			cells[y][x] = &Cell{
				Char: " ",
				FG:   "#ffffff",
				BG:   "",
			}
		}
	}

	return &Canvas{
		Width:     width,
		Height:    height,
		Cells:     cells,
		BGColor:   "",
		Name:      name,
		LatticeOn: false,
	}
}

func (c *Canvas) SetCell(x, y int, char string) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Cells[y][x].Char = char
	}
}

func (c *Canvas) GetCell(x, y int) *Cell {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		return c.Cells[y][x]
	}
	return nil
}

func (c *Canvas) Clear() {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.Cells[y][x].Char = " "
		}
	}
}

func GetComponentsDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "hang", "scab_components")
	os.MkdirAll(dir, 0755)
	return dir
}

func SaveCanvas(canvas *Canvas) error {
	dir := GetComponentsDir()
	filename := fmt.Sprintf("%s.json", canvas.Name)
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(canvas, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadCanvas(name string) (*Canvas, error) {
	dir := GetComponentsDir()
	filename := fmt.Sprintf("%s.json", name)
	path := filepath.Join(dir, filename)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var canvas Canvas
	err = json.Unmarshal(data, &canvas)
	if err != nil {
		return nil, err
	}

	canvas.BGColor = ""
	for _, row := range canvas.Cells {
		for _, cell := range row {
			if cell != nil && cell.BG == "#000000" {
				cell.BG = ""
			}
		}
	}

	return &canvas, nil
}

func ListComponents() []string {
	dir := GetComponentsDir()
	files, _ := filepath.Glob(filepath.Join(dir, "*.json"))

	var names []string
	for _, f := range files {
		name := filepath.Base(f)
		name = name[:len(name)-5]
		names = append(names, name)
	}
	return names
}

func (c *Canvas) Resize(newWidth, newHeight int) {
	newCells := make([][]*Cell, newHeight)
	for y := 0; y < newHeight; y++ {
		newCells[y] = make([]*Cell, newWidth)
		for x := 0; x < newWidth; x++ {
			if y < c.Height && x < c.Width {
				newCells[y][x] = c.Cells[y][x]
			} else {
				newCells[y][x] = &Cell{
					Char: " ",
					FG:   "#ffffff",
					BG:   "",
				}
			}
		}
	}

	c.Width = newWidth
	c.Height = newHeight
	c.Cells = newCells
}
