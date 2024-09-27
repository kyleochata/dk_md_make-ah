package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleochata/md_maker/badge"
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
	var cmd tea.Cmd
	if m.List, cmd = m.List.Update(msg); cmd != nil {
		return m, cmd
	}

	return m, nil
}
func (m Badge_model) View() string {
	return m.List.View()
}

func New_Badges_model(a Answers) tea.Model {
	items := PullData()
	li := make([]list.Item, len(items))
	for i, item := range items {
		li[i] = item
	}
	l := list.New(li, badge.CustomDelegate{}, a.Height, a.Width)
	return Badge_model{Answers: a, List: l}
}
