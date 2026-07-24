package shared

type Renderable struct {
	Buffer           [][]byte
	X, Y             int
	OriginX, OriginY int
	Transparent      bool
}
