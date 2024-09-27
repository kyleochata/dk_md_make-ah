package main

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/kyleochata/md_maker/badge"
)

type Intro_model struct {
	Answers
}

func (m Intro_model) Init() tea.Cmd { return nil }
func (m Intro_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m Intro_model) View() string {
	selectedBadges, ok := m.Responses["badge"].([]badge.Item)
	if !ok {
		return "Invalid badge data"
	}
	uiEl := []string{}
	for _, bi := range selectedBadges {
		// Process each badge item
		uiEl = append(uiEl, bi.Title())
	}
	return gloss.JoinVertical(gloss.Center, uiEl...)
}

func New_Intro_model(a Answers) tea.Model {
	return Intro_model{Answers: a}
}
