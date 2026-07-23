package ui

import (
	tea "github.com/charmbracelet/bubbletea"
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

func (m model) View() string {
	return "hello world!"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
