package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Use_model struct {
	Answers
	TextArea textarea.Model
}

const Usage string = "use"

func (m Use_model) Init() tea.Cmd { return nil }
func (m Use_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+b":
			return m.Send_to_Installation()
		case "ctrl+n":
			return m.send_to_license()
		}
	}
	var cmd tea.Cmd
	if m.TextArea, cmd = m.TextArea.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}
func (m Use_model) View() string {
	uiEl := []string{}
	uiEl = append(uiEl, m.titleText().Render("Please give instructions on how to use your project.\nUse Markdown syntax"))
	uiEl = append(uiEl, gloss.NewStyle().Align(gloss.Right).Render(m.TextArea.View()))
	uiEl = append(uiEl, "Ctrl+c: Quit | Ctrl+b: Previous Section | Ctrl+n: Next Section")
	return gloss.JoinVertical(gloss.Center, uiEl...)
}

// =============Helper========================
func New_Use_model(a Answers) tea.Model {
	ta := textarea.New()
	ta.SetHeight(a.Height / 2)
	ta.SetWidth(a.Width)
	if val, ok := a.Responses["use"].(string); ok {
		ta.SetValue(val)
	}
	ta.Focus()
	return Use_model{Answers: a, TextArea: ta}
}
func (m *Use_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}
func (m *Use_model) Send_to_Installation() (tea.Model, tea.Cmd) {
	m.save_use_data()
	return New_Install_model(m.Answers), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: m.Height,
			Width:  m.Width,
		}
	}
}
func (m Use_model) send_to_license() (tea.Model, tea.Cmd) {
	m.save_use_data()
	//check if we have already visited the license model and need to make a LICENSE file
	make, _ := m.Responses[License_w].(bool)
	if LicenseFileExists() {
		return New_has_License_model(m.Answers, FindLicenseType(), make), SendWindowMsg(m.Height, m.Width)
	} else {
		return New_Fail_License_check_model(m.Answers), func() tea.Msg {
			return tea.WindowSizeMsg{
				Height: m.Height,
				Width:  m.Width,
			}
		}
	}
}
func (m *Use_model) save_use_data() {
	m.Responses[Usage] = m.TextArea.Value()
}
func (m Use_model) titleText() gloss.Style {
	// Define the style here
	return gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("12")).
		Padding(1, 1, 1, 1).
		Margin(1, 1, 1, 1).
		Align(gloss.Center).
		Width(m.Width).
		Height(m.Height / 4)
}
