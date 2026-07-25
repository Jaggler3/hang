package sprite

import "hang.sh/internal/shared"

type Frame struct {
	Cells [][]*shared.Cell
}

type Sprite struct {
	Frames      []*Frame
	WalkCycleUD []int
	WalkCycleLR []int
	FrameRate   int
}

func (s *Sprite) WalkFrame(direction shared.Direction, animFrame int) int {
	switch direction {
	case shared.Left, shared.Right:
		return s.WalkCycleLR[animFrame%len(s.WalkCycleLR)]
	default:
		return s.WalkCycleUD[animFrame%len(s.WalkCycleUD)]
	}
}

func (s *Sprite) WalkCycleLen(direction shared.Direction) int {
	switch direction {
	case shared.Left, shared.Right:
		return len(s.WalkCycleLR)
	default:
		return len(s.WalkCycleUD)
	}
}

var frameTemplates = [][]string{
	{ // f0: idle
		` o `,
		`/|\`,
		`/ \`},
	{ // f1: walk sideways blink
		` o `,
		`/|\`,
		` |`},
	{ // f2: walk vertical left blink
		` o `,
		`/|\`,
		`' |`,
	},
	{ // f3: walk vertical left blink
		` o `,
		`/|\`,
		`| '`,
	},
}

func NewPlayerSprite(headChar, headFG, headBG string) *Sprite {
	frames := make([]*Frame, len(frameTemplates))
	for i, rows := range frameTemplates {
		frames[i] = frameFromRows(rows, headChar, headFG, headBG)
	}

	return &Sprite{
		Frames:      []*Frame{frames[0], frames[1], frames[2], frames[3]},
		WalkCycleUD: []int{0, 2, 0, 3},
		WalkCycleLR: []int{0, 1, 0, 1},
		FrameRate:   3,
	}
}

func frameFromRows(rows []string, headChar, headFG, headBG string) *Frame {
	cells := make([][]*shared.Cell, len(rows))
	for y, line := range rows {
		runes := []rune(line)
		row := make([]*shared.Cell, len(runes))
		for x, r := range runes {
			ch := string(r)
			fg := "#eeeeee"
			bg := headBG
			switch ch {
			case "o":
				ch = headChar
				fg = headFG
			case " ":
				bg = ""
			}
			row[x] = &shared.Cell{Char: ch, FG: fg, BG: bg}
		}
		cells[y] = row
	}
	return &Frame{Cells: cells}
}
