package data

import (
	"encoding/json"
	"os"
	"path/filepath"

	"hang.sh/internal/world"
)

type ContextKey string

const PlayerKey ContextKey = "player"
const WorldKey ContextKey = "world"

type tileData struct {
	Char     string `json:"char"`
	FG       string `json:"fg"`
	BG       string `json:"bg"`
	Walkable bool   `json:"walkable"`
}

type worldData struct {
	Width, Height int
	Tiles         [][]tileData
}

func worldPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "hang", "world.json")
}

func SaveWorld(w *world.World) error {
	w.RLock()
	defer w.RUnlock()

	data := worldData{
		Width:  w.Width,
		Height: w.Height,
		Tiles:  make([][]tileData, w.Height),
	}
	for y := 0; y < w.Height; y++ {
		row := make([]tileData, w.Width)
		for x := 0; x < w.Width; x++ {
			t := w.TileAt(x, y)
			if t != nil {
				row[x] = tileData{Char: t.Char, FG: t.FG, BG: t.BG, Walkable: t.Walkable}
			}
		}
		data.Tiles[y] = row
	}

	dir := filepath.Dir(worldPath())
	os.MkdirAll(dir, 0755)

	f, err := os.Create(worldPath())
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(data)
}

func LoadWorld(w *world.World) error {
	f, err := os.Open(worldPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	var data worldData
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	if data.Width != w.Width || data.Height != w.Height {
		return nil
	}

	w.Lock()
	defer w.Unlock()

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			td := data.Tiles[y][x]
			w.Tiles[y][x] = &world.Tile{Char: td.Char, FG: td.FG, BG: td.BG, Walkable: td.Walkable}
		}
	}

	return nil
}
