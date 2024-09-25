package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Badges_model struct {
	Answers
	Cursor int
	Picked []string
	List   list.Model
}

// list_item fufills the list.Item interface from bubbles/list
type list_item struct {
	name string
}

func (li list_item) FilterValue() string { return li.name }
func (li list_item) Title() string       { return li.name }
func (li list_item) Description() string { return li.name }

func New_Badges_model(a Answers) tea.Model {
	l := list.New([]list.Item{list_item{name: "golang"}}, list.NewDefaultDelegate(), a.Width/2, a.Height/2)
	l.Title = "Select the badges you want to display"
	l.InfiniteScrolling = true
	l.FilterInput.Placeholder = "Golang"
	l.FilterInput.TextStyle = gloss.NewStyle().Foreground(gloss.Color("#FFBF00")).Italic(true)
	l.FilterInput.Prompt = "Search badges: "
	l.FilterInput.PromptStyle = gloss.NewStyle().Foreground(gloss.Color("#FFFFFF")).Bold(true).MarginRight(2)
	return Badges_model{Answers: Answers{Responses: a.Responses}, List: l, Cursor: 0}
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
	bm.List, cmd = bm.List.Update(msg)
	return bm, cmd
}

func (bm Badges_model) View() string {
	return bm.List.View()
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
