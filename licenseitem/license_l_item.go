package licenseitem

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Item struct {
	License string `json:"license"`
	Summary string `json:"description"`
}

func (i Item) Title() string       { return i.License }
func (i Item) Description() string { return i.Summary }
func (i Item) FilterValue() string { return i.License }
func New_list_item(license, summary string) Item {
	return Item{License: license, Summary: summary}
}

var (
	normalTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).MarginLeft(1)
	selectedTextStyle = lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color("205")).Bold(true)
)

type CustomDelegate struct{}

func (d CustomDelegate) Height() int                              { return 1 }
func (d CustomDelegate) Spacing() int                             { return 0 }
func (d CustomDelegate) Update(ms tea.Msg, m *list.Model) tea.Cmd { return nil }

// func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
// 	i, ok := item.(Item)
// 	if !ok {
// 		return
// 	}
// 	isSelected := index == m.Index()
// 	var cursor string
// 	if isSelected {
// 		cursor = selectedTextStyle.Render("\t--> ")
// 	} else {
// 		cursor = "\t    "
// 	}
// 	availableWidth := m.Width() - lipgloss.Width(cursor) - 4
// 	summary := lipgloss.Wrap(i.Description(), availableWidth)
// 	name := normalTextStyle.Render(i.Title())
// 	if isSelected {
// 		name = selectedTextStyle.Render(i.Title())
// 		summary = selectedTextStyle.Render(i.Description())
// 	}
// 	fmt.Fprintf(w, "%s\t%s\t%s", cursor, name, summary)
// }

func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(Item)
	if !ok {
		return
	}

	isSelected := index == m.Index()
	// Get the available width from the list model (excluding cursor space)
	availableWidth := m.Width() - 4

	// Wrap the description and title based on the available width using wordwrap
	wrappedSummary := wordwrap.String(i.Description(), availableWidth)
	wrappedName := wordwrap.String(i.Title(), availableWidth)

	if isSelected {
		wrappedName = selectedTextStyle.Underline(true).Render(wrappedName)
		wrappedSummary = selectedTextStyle.Render(wrappedSummary)
	} else {
		wrappedName = normalTextStyle.Underline(true).Render(wrappedName)
		wrappedSummary = normalTextStyle.Render(wrappedSummary)
	}

	// Render the item, placing the name on the first line and the summary below it
	fmt.Fprintf(w, "%s\n%s", wrappedName, wrappedSummary)
}
func GetAvailableLicenses() ([]Item, error) {
	//path is from the root of the main directory where this func gets called
	file, err := os.ReadFile("licenseitem/available_licenses.json")
	if err != nil {
		log.Printf("Failed to read the JSON file: %v\n", err)
		return nil, err
	}

	var licenses []Item
	err = json.Unmarshal(file, &licenses)
	if err != nil {
		log.Printf("Failed to unmarshal the JSON file: %v\n", err)
		return nil, err
	}

	return licenses, nil
}
