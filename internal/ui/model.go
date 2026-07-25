package ui

import "charm.land/bubbles/v2/textinput"

type ScreenType int

const (
	ScreenStart ScreenType = iota
	ScreenGame
)

type Screen struct {
	ScreenType ScreenType
	Elements   []any
	HeadIndex  int
}

var HeadOptions = []string{"@", "o", "O", "X", "#", "&", "%", "$", "♥", "★", "♪", "0", "A", "M", "Z"}

func NewStartScreen() Screen {
	ti := textinput.New()
	ti.Placeholder = "your name"
	ti.Focus()
	ti.CharLimit = 20
	ti.SetWidth(20)

	return Screen{
		ScreenType: ScreenStart,
		Elements:   []any{ti},
		HeadIndex:  0,
	}
}
