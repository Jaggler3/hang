package world

import (
	"sync"

	"hang.sh/internal/chat"
	"hang.sh/internal/entity"
)

type World struct {
	Width, Height  int
	Tiles          [][]*Tile
	Players        map[string]*entity.Player
	NPCs           []*entity.NPC
	SpawnX, SpawnY int
	Messages       *chat.Buffer
	mu             sync.RWMutex
}

func New(width, height int) *World {
	w := &World{
		Width:    width,
		Height:   height,
		Players:  make(map[string]*entity.Player),
		SpawnX:   width / 2,
		SpawnY:   height / 2,
		Messages: chat.NewBuffer(),
	}

	w.Tiles = make([][]*Tile, height)
	for y := range w.Tiles {
		w.Tiles[y] = make([]*Tile, width)
		for x := range w.Tiles[y] {
			w.Tiles[y][x] = defaultTile()
		}
	}

	return w
}

func NewDefaultWorld() *World {
	w := New(80, 60)
	w.SpawnX = 40
	w.SpawnY = 30

	// some walls around the spawn
	for x := 35; x <= 45; x++ {
		w.setTile(x, 25, &Tile{Char: "#", FG: "#888888", BG: "#333333", Walkable: false})
		w.setTile(x, 35, &Tile{Char: "#", FG: "#888888", BG: "#333333", Walkable: false})
	}
	for y := 25; y <= 35; y++ {
		w.setTile(35, y, &Tile{Char: "#", FG: "#888888", BG: "#333333", Walkable: false})
		// w.setTile(45, y, &Tile{Char: "#", FG: "#888888", BG: "#333333", Walkable: false})
	}

	// water pond
	for y := 10; y <= 16; y++ {
		for x := 10; x <= 16; x++ {
			dx := x - 13
			dy := y - 13
			if dx*dx+dy*dy <= 12 {
				w.setTile(x, y, &Tile{Char: "~", FG: "#0088ff", BG: "#000044", Walkable: false})
			}
		}
	}

	// path
	for x := 36; x <= 44; x++ {
		w.setTile(x, 30, &Tile{Char: ".", FG: "#884422", BG: "#1a1a2e", Walkable: true})
	}

	// trees near spawn
	w.placeTree(50, 25)
	w.placeTree(52, 27)
	w.placeTree(48, 28)
	w.placeTree(55, 22)
	w.placeTree(30, 28)
	w.placeTree(28, 32)
	w.placeTree(50, 33)
	w.placeTree(55, 35)

	// flowers
	for i := 0; i < 12; i++ {
		x := 20 + i*3
		y := 40 + (i % 4) * 2
		w.setTile(x, y, &Tile{Char: "*", FG: "#ff44aa", BG: "#1a1a2e", Walkable: true})
	}

	return w
}

func (w *World) placeTree(x, y int) {
	w.setTile(x, y, &Tile{Char: "A", FG: "#22cc44", BG: "#1a1a2e", Walkable: false})
	w.setTile(x+1, y, &Tile{Char: "A", FG: "#22cc44", BG: "#1a1a2e", Walkable: false})
	w.setTile(x, y-1, &Tile{Char: "A", FG: "#33ee55", BG: "#1a1a2e", Walkable: false})
	w.setTile(x+1, y-1, &Tile{Char: "A", FG: "#33ee55", BG: "#1a1a2e", Walkable: false})
	w.setTile(x, y+1, &Tile{Char: "|", FG: "#884422", BG: "#1a1a2e", Walkable: false})
	w.setTile(x+1, y+1, &Tile{Char: "|", FG: "#884422", BG: "#1a1a2e", Walkable: false})
}

func (w *World) AddPlayer(p *entity.Player) {
	w.mu.Lock()
	defer w.mu.Unlock()
	p.X = w.SpawnX
	p.Y = w.SpawnY
	w.Players[p.ID] = p
}

func (w *World) RemovePlayer(id string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.Players, id)
}

func (w *World) UpdatePlayerPosition(id string, x, y int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p, ok := w.Players[id]; ok {
		p.X = x
		p.Y = y
	}
}

func (w *World) ForEachPlayer(fn func(*entity.Player)) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, p := range w.Players {
		fn(p)
	}
}

func (w *World) ForEachPlayerBreak(fn func(*entity.Player) bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, p := range w.Players {
		if fn(p) {
			break
		}
	}
}

func (w *World) ForEachNPC(fn func(*entity.NPC)) {
	for _, n := range w.NPCs {
		fn(n)
	}
}

func (w *World) BroadcastChat(name, text string) {
	w.Messages.Add(name, text)
}

func (w *World) RLock() {
	w.mu.RLock()
}

func (w *World) RUnlock() {
	w.mu.RUnlock()
}

func (w *World) Lock() {
	w.mu.Lock()
}

func (w *World) Unlock() {
	w.mu.Unlock()
}

func (w *World) IsWalkable(x, y int) bool {
	if x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return false
	}
	return w.Tiles[y][x].Walkable
}

func (w *World) TileAt(x, y int) *Tile {
	if x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return nil
	}
	return w.Tiles[y][x]
}

func defaultTile() *Tile {
	return &Tile{Char: ".", FG: "#555555", BG: "#1a1a2e", Walkable: true}
}

func (w *World) setTile(x, y int, t *Tile) {
	if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
		w.Tiles[y][x] = t
	}
}
