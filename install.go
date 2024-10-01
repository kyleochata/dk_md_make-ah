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
		// return m, nil
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

func (m Installation_model) View() string {
	if len(m.List.Items()) == 0 {
		return "No items pop"
	}
	uiEl := []string{}
	if m.FocusState == fs_im_list {
		s := "Choose a pre-built installation guide for languages/tools used for your project.\nSelecting one will add generic install instructions for the language/tool picked."
		uiEl = append(uiEl, m.titleStyle().Render(s))
	}
	uiEl = append(uiEl, gloss.NewStyle().Margin(1, 0, 1, 0).Render(m.List.View()))
	if m.FocusState == fs_im_ta {
		uiEl = append(uiEl, m.ta_prompt_style())
	}
	uiEl = append(uiEl, m.TextArea.View())
	if m.FocusState == fs_im_ta {
		s := "Ctrl+C: Quit | Ctrl+N: Next Section | Tab: Switch to List"
		uiEl = append(uiEl, m.titleStyle().Render(s))
	}

	return gloss.JoinVertical(gloss.Center, uiEl...)
}

// =========Helper==================
var available_installs = []string{"C", "C++", "C#", "Golang", "Node.js", "Rails", "Ruby", "TypeScript"}

func (m *Installation_model) windowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	m.List.SetHeight(m.Height / 3)
	m.List.SetWidth(m.Width) // Ensure full width
	m.TextArea.SetHeight(m.Height / 2)
	m.TextArea.SetWidth(m.Width) // Ensure full width
}

func New_Install_model(a Answers) tea.Model {
	list_items := make_list_items(available_installs)
	list := list.New(list_items, install.CustomDelegate{}, a.Width-3, a.Height/3)
	list.Title = "Available Install Boilerplate"
	ta := textarea.New()
	ta.SetWidth(a.Width - 3)
	ta.SetHeight(a.Height / 2)
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
	} else {
		m.FocusState = fs_im_list
		m.List.SetShowHelp(true)
		m.TextArea.Blur()
	}
}

// ===========Style==========================
func (m Installation_model) titleStyle() gloss.Style {
	return gloss.NewStyle().
		Foreground(gloss.Color("#FFBF00")).Align(gloss.Right).Margin(1, 1, 1, 1)
}

func (m Installation_model) ta_prompt_style() string {
	blueTextStyle := gloss.NewStyle().
		Foreground(gloss.Color("12"))
	boldYellowTextStyle := gloss.NewStyle().
		Bold(true).
		Foreground(gloss.Color("#FFBF00"))

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
