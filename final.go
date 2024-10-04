package main

import tea "github.com/charmbracelet/bubbletea"

type final_model struct {
	Answers
}

func (m final_model) Init() tea.Cmd { return nil }
func (m final_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m final_model) View() string {
	return "from final"
}
func (m *final_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}
func New_final_model(a Answers) tea.Model {
	return final_model{Answers: a}
}
