package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Badge_model struct {
	Answers
	List list.Model
}

func (m Badge_model) Init() tea.Cmd { return nil }
func (m Badge_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width
		m.List.SetHeight(m.Height)
		m.List.SetWidth(m.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m Badge_model) View() string {
	return m.List.View()
}
