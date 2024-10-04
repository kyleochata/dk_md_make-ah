package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	Wild string = "wild"
)

type Wild_model struct {
	Answers
	textarea textarea.Model
}

func (m Wild_model) Init() tea.Cmd { return nil }
func (m Wild_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+b":
			m.save_wild_data()
			return Send_to_Contributors(m.Answers)
		case "ctrl+n":
			m.save_wild_data()
			return Send_to_final(m.Answers)
		}
	}
	var cmd tea.Cmd
	if m.textarea, cmd = m.textarea.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}
func (m Wild_model) View() string {
	return m.textarea.View()
}
func New_wild_model(a Answers) tea.Model {
	ta := textarea.New()
	ta.SetWidth(a.Width)
	s_val, ok := a.Responses[Wild].(string)
	if ok {
		ta.SetValue(s_val)
	}
	ta.Focus()
	return Wild_model{Answers: a, textarea: ta}
}
func (m *Wild_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	m.textarea.SetWidth(m.Width)
}
func (m *Wild_model) save_wild_data() {
	m.Responses[Wild] = m.textarea.Value()
}
func Send_to_final(a Answers) (tea.Model, tea.Cmd) {
	return New_final_model(a), SendWindowMsg(a.Height, a.Width)
}
