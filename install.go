package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleochata/md_maker/install"
)

const (
	fs_im_list fs_im = iota
	fs_im_ta
)

type fs_im int
type Installation_model struct {
	Answers
	List           list.Model
	TextArea       textarea.Model
	FocusState     fs_im
	InstallChoices []string
}

func (m Installation_model) Init() tea.Cmd { return nil }
func (m Installation_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m Installation_model) View() string {
	if len(m.List.Items()) == 0 {
		return "No items pop"
	}
	return m.List.View()
}

// =========Helper==================
var available_installs = []string{"C", "C++", "C#", "Golang", "Node.js", "Rails", "Ruby", "TypeScript"}

func New_Install_model(a Answers) tea.Model {
	list_items := make_list_items(available_installs)
	list := list.New(list_items, install.CustomDelegate{}, a.Width, a.Height/3)
	list.Title = "Available Install Boilerplate"
	ta := textarea.New()
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
