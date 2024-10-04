package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Answers struct {
	Responses     map[string]any
	Height, Width int
}
type Title_model struct {
	Answers
	TextInput  textinput.Model
	FocusState FocusState
}
type FocusState int

const (
	Title            string     = "title"
	FocusState_Input FocusState = iota
	FocusState_Title
)

func (t Title_model) Init() tea.Cmd { return nil }
func (t Title_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.Height, t.Width = msg.Height, msg.Width
		return t, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			t.toggleFocus()
			return t, nil
		case "ctrl+c", "esc":
			return t, tea.Quit
		case "enter":
			if t.FocusState == FocusState_Input {
				return t.Send_to_Badges()
			}
		}
		var cmd tea.Cmd
		if t.TextInput, cmd = t.TextInput.Update(msg); cmd != nil {
			return t, cmd
		}
		return t, nil
	}

	return t, nil
}
func (t Title_model) View() string {
	uiEl := []string{t.titleView()}
	uiEl = append(uiEl, t.textInputView())
	uiEl = append(uiEl, "Press Enter to progress")
	uiEl = append(uiEl, "Ctrl+C or Esc: Quit")

	return gloss.JoinVertical(gloss.Center, uiEl...)
}

// =============helper================
func New_Title_model(a Answers) tea.Model {
	ti := textinput.New()
	//First model after open model. Need to init the response map if first visit to title model
	if a.Responses == nil {
		a.Responses = map[string]any{Title: nil}
	}
	if a.Responses[Title] != nil {
		ti.SetValue(a.Responses[Title].(string))
	}
	ti.Focus()
	return Title_model{Answers: a, TextInput: ti, FocusState: FocusState_Input}
}

func (t Title_model) Send_to_Badges() (tea.Model, tea.Cmd) {
	t.Responses[Title] = t.TextInput.Value()
	return New_Badges_model(t.Answers), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: t.Height,
			Width:  t.Width,
		}
	}
}

func (t Title_model) toggleFocus() {
	if t.FocusState == FocusState_Input {
		t.FocusState = FocusState_Title
		t.TextInput.Blur()
	} else {
		t.FocusState = FocusState_Input
		t.TextInput.Focus()
	}
}

func (t Title_model) titleView() string {
	s := "What would you like as the title of this README.md?"
	if t.FocusState == FocusState_Title {
		return t.titleTextFocus().Render(s)
	}
	return t.titleText().Render(s)
}
func (t Title_model) textInputView() string {
	var s strings.Builder
	s.WriteString(t.TextInput.View())
	if t.FocusState == FocusState_Input {
		return t.textInputFocus().Render(s.String())
	} else {
		return t.textInput().Render(s.String())
	}
}

// ============style===================
func (t Title_model) titleTextFocus() gloss.Style {
	// Define the style here
	return gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("12")).
		Padding(1, 1, 1, 1).
		Margin(1, 1, 1, 1).
		Align(gloss.Center).
		Width(t.Width).
		Height(t.Height / 4)
}
func (t Title_model) titleText() gloss.Style {
	return gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("12")).
		Padding(1, 1, 1, 1).
		Margin(1, 1, 1, 1).
		Align(gloss.Center).
		Width(t.Width).
		Height(t.Height / 4)
}
func (t Title_model) textInputFocus() gloss.Style {
	return gloss.NewStyle().
		Foreground(gloss.Color("FFFFFF")).
		Padding(1, 1, 1, 1).
		Margin(1, 1, 1, 1).
		Bold(true).
		Width(t.Width / 2).
		Height(t.Height / 4)
}
func (t Title_model) textInput() gloss.Style {
	return gloss.NewStyle().
		Foreground(gloss.Color("B2BEB5")).
		Padding(1, 1, 1, 1).
		Margin(1, 1, 1, 1).
		Align(gloss.Center).
		Width(t.Width / 2).
		Height(t.Height / 4).
		Italic(true)
}
