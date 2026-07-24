package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"
	"hang.sh/internal/data"
	"hang.sh/internal/player"
)

type model struct {
	width, height int
	player        *player.Player
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() tea.View {
	view := tea.NewView("hello " + m.player.Name)
	view.AltScreen = true
	return view
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func TeaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	ctx := sess.Context()
	p := ctx.Value(data.PlayerKey).(*player.Player)
	return model{player: p}, []tea.ProgramOption{}
}
