package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
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
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.HandleSelectBadge()
			// Ensure the list height is updated after selection
			lh := m.Height * 2 / 3
			m.List.SetHeight(lh)
			return m, nil

		case "ctrl+n":
			return m.Send_to_Intro(m.Answers)
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
		badgeEl = append(badgeEl, gloss.NewStyle().Render(s))
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
	l.AdditionalFullHelpKeys = AdditionalFullHelpKeys
	l.AdditionalShortHelpKeys = AdditionalShortHelpKeys
	l.AdditionalShortHelpKeys()
	return Badge_model{Answers: a, List: l, BadgeChoices: badges_arr}
}

func (m *Badge_model) HandleSelectBadge() {
	selectedItem := m.List.SelectedItem()
	if selectedItem == nil {
		return
	}
	badgeItem, ok := selectedItem.(badge.Item)
	if !ok || badgeItem.IsSection {
		return
	}
	// Toggle the badge selection
	for i, item := range m.BadgeChoices {
		if item.Name == badgeItem.Name {
			m.BadgeChoices = append(m.BadgeChoices[:i], m.BadgeChoices[i+1:]...)
			badgeItem.BadgePicked = false
			return
		}
	}

	// Add the badge if not already selected
	badgeItem.BadgePicked = true
	m.BadgeChoices = append(m.BadgeChoices, badgeItem)
}
func (m *Badge_model) TurnOffHelp() {
	m.List.SetShowHelp(false)
}

// CHORE: Fix the args
func (m *Badge_model) Send_to_Intro(a Answers) (tea.Model, tea.Cmd) {
	m.Responses["badge"] = m.BadgeChoices
	m.TurnOffHelp()
	return New_Intro_model(m.Answers), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: m.Height,
			Width:  m.Width,
		}
	}
}

// Short help keys
func AdditionalShortHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "Next Section"),
		),
	}
}

// Full help keys
func AdditionalFullHelpKeys() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "Next Section"),
		),
		key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select Badge"),
		),
		key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Quit"),
		),
	}
}
