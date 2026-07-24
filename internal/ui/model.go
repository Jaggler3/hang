package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/ssh"
)

type model struct {
}

func initialModel() *model {
	return &model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() tea.View {
	view := tea.NewView("hello world!")
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
	}
	return m, nil
}

func TeaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	return initialModel(), []tea.ProgramOption{}
}
