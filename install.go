package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/kyleochata/md_maker/install"
)

const (
	fs_im_list fs_im = iota
	fs_im_ta
)

type fs_im int
type Installation_model struct {
	Answers
	//InstallChoices []string
	List           list.Model
	TextArea       textarea.Model
	FocusState     fs_im
	InstallChoices []string
}

func (m Installation_model) Init() tea.Cmd { return nil }
func (m Installation_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowResize(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.toggle_focus()
			return m, nil
		}
	}
	var cmd tea.Cmd
	if m.FocusState == fs_im_list {
		if m.List, cmd = m.List.Update(msg); cmd != nil {
			return m, cmd
		}
	} else {
		if m.TextArea, cmd = m.TextArea.Update(msg); cmd != nil {
			return m, cmd
		}
	}
	return m, nil
}

// func (m Installation_model) View() string {
// 	if len(m.List.Items()) == 0 {
// 		return "No items pop"
// 	}
// 	uiEl := []string{}
// 	if m.FocusState == fs_im_list {
// 		s := "Choose a pre-built installation guide for languages/tools used for your project.\nSelecting one will add generic install instructions for the language/tool picked."
// 		uiEl = append(uiEl, m.titleStyle().Width(m.Width).Render(s))
// 	}
// 	uiEl = append(uiEl, gloss.NewStyle().Margin(1, 0, 1, 0).Render(m.List.View()))
// 	if m.FocusState == fs_im_ta {
// 		uiEl = append(uiEl, m.ta_prompt_style())
// 	}
// 	uiEl = append(uiEl, m.TextArea.View())
// 	if m.FocusState == fs_im_ta {
// 		s := "Ctrl+C: Quit | Ctrl+N: Next Section | Tab: Switch to List"
// 		uiEl = append(uiEl, m.titleStyle().Render(s))
// 	}

//		return gloss.JoinVertical(gloss.Center, uiEl...)
//	}
func (m Installation_model) View() string {
	if len(m.List.Items()) == 0 {
		return "No items pop"
	}
	uiEl := []string{}

	// Title or instructions centered
	if m.FocusState == fs_im_list {
		s := "Choose a pre-built installation guide for languages/tools used for your project.\nSelecting one will add generic install instructions for the language/tool picked."
		uiEl = append(uiEl, m.titleStyle().Render(s))
	}

	// Render the list
	uiEl = append(uiEl, gloss.NewStyle().Align(gloss.Center).Margin(1, 0, 1, 0).Render(m.List.View()))

	// Render the textarea and prompt if focused
	if m.FocusState == fs_im_ta {
		uiEl = append(uiEl, m.ta_prompt_style()) // Add prompt only when focused
	}
	uiEl = append(uiEl, gloss.NewStyle().Align(gloss.Center).Padding(1, 2, 1, 2).Margin(1, 0, 1, 0).Render(m.TextArea.View()))

	// Add footer help information if focused on textarea
	if m.FocusState == fs_im_ta {
		s := "Ctrl+C: Quit | Ctrl+N: Next Section | Tab: Switch to List"
		uiEl = append(uiEl, m.titleStyle().Align(gloss.Center).Render(s))
	}

	// Use JoinVertical with added styles for centering
	return gloss.JoinVertical(gloss.Center, uiEl...)
}

// =========Helper==================
var available_installs = []string{"C", "C++", "C#", "Golang", "Node.js", "Rails", "Ruby", "TypeScript"}

func (m *Installation_model) windowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	m.List.SetHeight(m.Height / 3)
	m.List.SetWidth(m.Width)
	m.TextArea.SetWidth(m.Width) // Set textarea width to 2/3 of total width
}

func New_Install_model(a Answers) tea.Model {
	list_items := make_list_items(available_installs)
	list := list.New(list_items, install.CustomDelegate{}, a.Width, a.Height/3)
	list.Title = "Available Install Boilerplate"
	ta := textarea.New()
	ta.SetHeight(a.Height / 2)
	ta.SetWidth(a.Width)
	ta.Blur()
	return Installation_model{Answers: a, List: list, TextArea: ta, FocusState: fs_im_list, InstallChoices: []string{}}
}

func make_list_items(xs []string) []list.Item {
	var items []install.Item
	for _, item := range xs {
		items = append(items, install.New_list_item(item))
	}
	li := make([]list.Item, len(items))
	for i, item := range items {
		li[i] = item
	}
	return li
}

func (m *Installation_model) toggle_focus() {
	if m.FocusState == fs_im_list {
		m.List.SetShowHelp(false)
		m.FocusState = fs_im_ta
		m.TextArea.Focus()
		m.TextArea.Prompt = m.ta_prompt_style() // Set prompt when gaining focus
	} else {
		m.FocusState = fs_im_list
		m.List.SetShowHelp(true)
		m.TextArea.Blur()
		m.TextArea.Prompt = "" // Clear prompt when losing focus
	}
}

// ===========Style==========================
func (m Installation_model) titleStyle() gloss.Style {
	return gloss.NewStyle().
		Foreground(gloss.Color("#FFBF00")).
		Align(gloss.Center).
		Margin(1, 1, 1, 1)
}

func (m Installation_model) ta_prompt_style() string {
	blueTextStyle := gloss.NewStyle().
		Foreground(gloss.Color("12")).
		Align(gloss.Center)
	boldYellowTextStyle := gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("#FFBF00")).
		Align(gloss.Center)

	introText := "If your project requires further installation steps in addition to the pre-generated install guides, please add them below. Use "
	use := "Use "
	markdownText := "markdown syntax"
	endingText := " for best results"

	return fmt.Sprintf("%s\n%s%s%s",
		blueTextStyle.Render(introText),
		blueTextStyle.Render(use),
		boldYellowTextStyle.Render(markdownText),
		blueTextStyle.Render(endingText),
	)
}
