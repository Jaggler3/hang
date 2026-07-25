package shared

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Cell struct {
	Char      string `json:"char"`
	FG        string `json:"fg"`
	BG        string `json:"bg"`
	Bold      bool   `json:"bold"`
	Italic    bool   `json:"italic"`
	Underline bool   `json:"underline"`
}

type Renderable struct {
	Cells       [][]*Cell
	X, Y        int
	OriginX, OriginY int
	Transparent bool
}
