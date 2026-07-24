package player

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Player struct {
	X, Y      int
	Direction Direction
	Name      string
}

func NewPlayer() *Player {
	return &Player{
		X:         0,
		Y:         0,
		Direction: Down,
		Name:      "Martin",
	}
}
