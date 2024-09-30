package badge

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Name        string
	Badge       string
	IsSection   bool
	BadgePicked bool //If item is in the selected pool of badges

}

func (i Item) Title() string {
	return i.Name
}
func (i Item) Description() string {
	if i.IsSection {
		return ""
	}
	return i.Badge
}
func (i Item) FilterValue() string { return i.Name }
func NewItem(name, badge string) Item {
	return Item{Name: name, Badge: badge}
}

var (
	normalTextStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalBadgeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	selectedBadgeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	badgePickedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Background(lipgloss.Color("252"))
)

type CustomDelegate struct{}

func (d CustomDelegate) Height() int {
	return 1
}
func (d CustomDelegate) Spacing() int {
	return 0
}
func (d CustomDelegate) Update(ms tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(Item) // Use the correct cast to badge.Item
	if !ok {
		return
	}

	// Check if the current item is selected
	isSelected := index == m.Index()

	// Selection marker
	var cursor string
	if isSelected {
		cursor = "--> " // Customize your cursor here
	} else {
		cursor = "    " // Empty when not selected
	}
	// Apply styles based on whether the item is selected or not
	name := normalTextStyle.Render(i.Name)
	badge := normalBadgeStyle.Render(i.Badge)
	if i.BadgePicked {
		name = badgePickedStyle.Render(i.Name)
		badge = badgePickedStyle.Render(i.Badge)
	}

	if isSelected {
		name = selectedTextStyle.Render(i.Name)
		badge = selectedBadgeStyle.Render(i.Badge)
	}

	// Render the combined item output
	fmt.Fprintf(w, "%s  %s\t%s", cursor, name, badge) // Render name and badge with a tab for spacing
}
