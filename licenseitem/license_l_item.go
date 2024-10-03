package licenseitem

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Name string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return "" }
func (i Item) FilterValue() string { return i.Name }
func New_list_item(s string) Item {
	return Item{Name: s}
}

var (
	normalTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	selectedTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
)

type CustomDelegate struct{}

func (d CustomDelegate) Height() int                              { return 1 }
func (d CustomDelegate) Spacing() int                             { return 0 }
func (d CustomDelegate) Update(ms tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(Item)
	if !ok {
		return
	}
	isSelected := index == m.Index()
	var cursor string
	if isSelected {
		cursor = selectedTextStyle.Render("\t--> ")
	} else {
		cursor = "\t    "
	}
	name := normalTextStyle.Render(i.Name)
	if isSelected {
		name = selectedTextStyle.Render(i.Name)
	}
	fmt.Fprintf(w, "%s\t%s", cursor, name)
}
