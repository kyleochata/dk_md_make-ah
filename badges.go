package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Badges_model struct {
	Answers
	Cursor   int
	Picked   []string
	TextArea textarea.Model
}

func New_Badges_model(a Answers) tea.Model {
	ta := textarea.New()
	ta.Placeholder = "This project was built using...to solve...."
	ta.Focus()
	ta.Prompt = "Please introduce this project to explain the why and how of the project."

	ta.FocusedStyle.Prompt.Foreground(gloss.Color("FFA500")).Bold(true)

	return Badges_model{Answers: Answers{Responses: a.Responses}, TextArea: ta, Cursor: 0}
}

func (bm Badges_model) Init() tea.Cmd {
	return nil
}

func (bm Badges_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		bm.Height, bm.Width = msg.Height, msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return bm, tea.Quit
		// case "enter":
		case "left":
			return bm.Send_to_title()
		}
	}
	var cmd tea.Cmd
	bm.TextArea, cmd = bm.TextArea.Update(msg)
	return bm, cmd
}

func (bm Badges_model) View() string {
	return bm.TextArea.View()
}

// ===========================helper============================
func (bm Badges_model) Send_to_title() (tea.Model, tea.Cmd) {
	//doesn't matter if empty. Init this key-val so that we can return here without hitting "enter" in title view
	bm.Responses["badges"] = true
	return New_Title_model(bm.Answers), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: bm.Height,
			Width:  bm.Width,
		}
	}
}

//========================style================================
// func (bm Badges_model)
