package ui

import (
	"bytes"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"
	"hang.sh/internal/data"
	"hang.sh/internal/player"
)

type model struct {
	width, height int
	sized         bool
	player        *player.Player
	buffer        [][]byte
}

func (m model) Init() tea.Cmd {
	return nil
}

func setupView(view *tea.View) {
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion
}

func (m model) View() tea.View {
	if !m.sized {
		// first frame for getting screen size
		view := tea.NewView("")
		setupView(&view)
		return view
	}

	// layers:
	// - world: ground, sprites, players
	// - billboards
	// - ui: buttons, cards, text
	// - popups
	//

	// clear screen
	m.buffer = make([][]byte, m.height)
	for i := range m.buffer {
		m.buffer[i] = bytes.Repeat([]byte{' '}, m.width)
	}

	// WORLD
	worldOffsetX, worldOffsetY := m.width/2, m.height/2
	playerElement := [][]byte{
		[]byte(" o "),
		[]byte("/|\\"),
		[]byte(" n "),
	}
	m.renderElement(playerElement, worldOffsetX, worldOffsetY, 1, 1, false)

	view := tea.NewView(m.formatBuffer())
	setupView(&view)
	return view
}

func (m model) formatBuffer() string {
	lines := make([]string, len(m.buffer))
	for i, row := range m.buffer {
		lines[i] = string(row)
	}

	return strings.Join(lines, "\n")
}

func (m *model) renderCell(fill byte, position_x int, position_y int, erase bool) {
	size_y := len(m.buffer)
	if size_y == 0 || position_y < 0 || position_y > size_y {
		return
	}
	size_x := len(m.buffer[0])
	if size_x == 0 || position_x < 0 || position_x > size_x {
		return
	}
	if fill == ' ' && !erase {
		return
	}
	m.buffer[position_y][position_x] = fill
}

func (m *model) renderElement(fill [][]byte, position_x int, position_y int, origin_x int, origin_y int, erase bool) {
	size_y := len(fill)
	if size_y == 0 {
		return
	}
	size_x := len(fill[0])
	if size_x == 0 {
		return
	}
	for y := range size_y {
		for x := range size_x {
			m.renderCell(fill[y][x], position_x+x-origin_x, position_y+y-origin_y, erase)
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg: // first frame
		m.width = msg.Width
		m.height = msg.Height
		m.sized = true
	}
	return m, nil
}

func TeaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	ctx := sess.Context()
	p := ctx.Value(data.PlayerKey).(*player.Player)
	return model{player: p}, []tea.ProgramOption{}
}
