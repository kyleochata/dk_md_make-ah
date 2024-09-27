package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/kyleochata/md_maker/badge"
)

type Badge_model struct {
	Answers
	List         list.Model
	BadgeChoices []badge.Item
}

func (m Badge_model) Init() tea.Cmd { return nil }
func (m Badge_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height, m.Width = msg.Height, msg.Width
		lh := m.Height * 2 / 3
		m.List.SetHeight(lh)
		m.List.SetWidth(m.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.HandleSelectBadge()
			return m, nil
		}
	}
	var cmd tea.Cmd
	if m.List, cmd = m.List.Update(msg); cmd != nil {
		return m, cmd
	}

	return m, nil
}
func (m Badge_model) View() string {
	uiEl := []string{"Choose badges"}
	badgeEl := []string{}
	for _, chosenBadge := range m.BadgeChoices {
		s := fmt.Sprintf("%s\t", chosenBadge.Name)
		badgeEl = append(badgeEl, s)
	}
	uiEl = append(uiEl, gloss.JoinHorizontal(gloss.Center, badgeEl...))
	uiEl = append(uiEl, m.List.View())
	return gloss.JoinVertical(gloss.Center, uiEl...)
}

//===============Helper===========================

func New_Badges_model(a Answers) tea.Model {
	items := PullData()
	li := make([]list.Item, len(items))
	for i, item := range items {
		li[i] = item
	}
	badges_picked := a.Responses["badge"]
	var badges_arr []badge.Item
	if badges_picked == nil {
		badges_picked = []badge.Item{}
	} else {
		//type assertion from any to []badge.Item
		if bp, ok := badges_picked.([]badge.Item); ok {
			badges_arr = append(badges_arr, bp...)
		}
	}
	lh := a.Height * 2 / 3
	l := list.New(li, badge.CustomDelegate{}, a.Width, lh)
	return Badge_model{Answers: a, List: l, BadgeChoices: badges_arr}
}

func (m *Badge_model) HandleSelectBadge() {
	selectedItem := m.List.SelectedItem()
	if selectedItem == nil {
		return // No item selected
	}
	badgeItem, ok := selectedItem.(badge.Item)
	if !ok || badgeItem.IsSection {
		return // Ensure the type is correct
	}
	// Toggle the badge selection
	for i, item := range m.BadgeChoices {
		if item.Name == badgeItem.Name {
			// Remove the badge if already selected
			m.BadgeChoices = append(m.BadgeChoices[:i], m.BadgeChoices[i+1:]...)
			badgeItem.BadgePicked = false // Update the badgeItem
			return
		}
	}

	// Add the badge if not already selected
	badgeItem.BadgePicked = true // Update the badgeItem
	m.BadgeChoices = append(m.BadgeChoices, badgeItem)
}
