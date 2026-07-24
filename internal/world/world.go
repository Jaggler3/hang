package world

import "hang.sh/internal/shared"

// we can plug in some params here so we can only show local elements etc eventually
// just cant pass the whole render model or we end up with an import loop
func Elements() []shared.Renderable {
	currentPlayer := shared.Renderable{
		Buffer: [][]byte{
			[]byte(" o "),
			[]byte("/|\\"),
			[]byte(" n "),
		},
		OriginX: 1,
		OriginY: 1,
	}

	elements := []shared.Renderable{currentPlayer}
	return elements
}
