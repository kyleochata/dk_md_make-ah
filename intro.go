package main

import (
	ta "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type Intro_model struct {
	Answers
	TextArea ta.Model
}

func (m Intro_model) Init() tea.Cmd { return nil }
func (m Intro_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	}
	var cmd tea.Cmd
	if m.TextArea, cmd = m.TextArea.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}
func (m Intro_model) View() string {
	return m.TextArea.View()
}

func New_Intro_model(a Answers) tea.Model {
	ta := ta.New()
	ta.Focus()
	return Intro_model{Answers: a, TextArea: ta}
}

//Plop into view to see data drilling
// selectedBadges, ok := m.Responses["badge"].([]badge.Item)
// if !ok {
// 	return "Invalid badge data"
// }
// uiEl := []string{}
// for _, bi := range selectedBadges {
// 	// Process each badge item
// 	uiEl = append(uiEl, bi.Title())
// }
// return gloss.JoinVertical(gloss.Center, uiEl...)
