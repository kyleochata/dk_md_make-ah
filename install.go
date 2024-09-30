package main

import tea "github.com/charmbracelet/bubbletea"

type Installation_model struct {
	Answers
}

func (m Installation_model) Init() tea.Cmd { return nil }
func (m Installation_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m Installation_model) View() string {
	return "hello"
}
