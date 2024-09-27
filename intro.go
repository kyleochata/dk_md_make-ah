package main

import (
	ta "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
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
		case "ctrl+n":
			return Send_to_Installation(m)
		}

	}
	var cmd tea.Cmd
	if m.TextArea, cmd = m.TextArea.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}
func (m Intro_model) View() string {
	uiEl := []string{}
	uiEl = append(uiEl, m.TitleRender())
	uiEl = append(uiEl, m.TextArea.View())
	return gloss.JoinVertical(gloss.Center, uiEl...)
}

func New_Intro_model(a Answers) tea.Model {
	ta := ta.New()
	ta.SetWidth(a.Width)
	ta.Focus()
	return Intro_model{Answers: a, TextArea: ta}
}
func (m Intro_model) TitleRender() string {
	return m.TitleStyle().Render("In a few sentences, describe your project. What problem does this project solve?")
}
func GetTextAreaValue(m Intro_model) string {
	return m.TextArea.Value()
}
func Send_to_Installation(m Intro_model) (tea.Model, tea.Cmd) {
	m.Responses["intro"] = GetTextAreaValue(m)
	return Installation_model{Answers: m.Answers}, func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: m.Height,
			Width:  m.Width,
		}
	}
}

type Installation_model struct {
	Answers
}

func (m Installation_model) Init() tea.Cmd { return nil }
func (m Installation_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
func (m Installation_model) View() string {
	return "hello"
}

// =======Style==================
func (m Intro_model) TitleStyle() gloss.Style {
	return gloss.NewStyle().
		Foreground(gloss.Color("#FFBF00")).
		Align(gloss.Right).
		Margin(1, 1, 1, 1)
}
