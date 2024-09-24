package main

import (
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type OpenModel struct {
	Height, Width int
	TimeUp        bool
	Count         int
}

func InitialModel() OpenModel {
	return OpenModel{Count: 3}
}

type TickMsg time.Time

func (om OpenModel) Init() tea.Cmd {
	return om.doTick()
}
func (om OpenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		if om.Count > 0 {
			om.Count--
			if om.Count == 0 {
				return om, tea.Quit
			} else {
				return om, om.doTick()
			}
		}
	case tea.WindowSizeMsg:
		om.Height, om.Width = msg.Height, msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return om, tea.Quit
		case "enter":
			return Send_from_open(om)
		}
	}
	return om, nil
}

// not getting called after every TickMsg received
func (om OpenModel) View() string {
	s := strconv.Itoa(om.Count)
	return s
}

// ======================helper===========================================

func (om OpenModel) doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
func Send_from_open(om OpenModel) (tea.Model, tea.Cmd) {
	return New_Title_model(Answers{}), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: om.Height,
			Width:  om.Width,
		}
	}
}

//======================style===============================================
