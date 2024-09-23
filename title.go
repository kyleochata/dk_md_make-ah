package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Title struct {
	Answers       map[string]any
	Height, Width int
	TitleInput    textarea.Model
}

func (t Title) Init() tea.Cmd { return nil }
func (t Title) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.Height, t.Width = msg.Height, msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return t, tea.Quit
		}
	}
	return t, nil
}
func (t Title) View() string {
	return "title screen"
}
