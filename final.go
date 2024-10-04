package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleochata/md_maker/badge"
)

type final_model struct {
	Answers
	content string
}

func (m final_model) Init() tea.Cmd { return nil }
func (m final_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m final_model) View() string {
	m.compileUserResponse()
	return "from final"
}
func (m *final_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}
func New_final_model(a Answers) tea.Model {
	return final_model{Answers: a}
}
func (m *final_model) compileUserResponse() {
	var b strings.Builder
	//title
	s_title := m.Responses[Title].(string)
	b.WriteString(fmt.Sprintf("# %s\n", s_title))
	b.WriteString("---\n")
	//badge
	xbitem_badges := m.Responses[Badge].([]badge.Item)
	for i, badge := range xbitem_badges {
		if i == len(xbitem_badges)-1 {

			b.WriteString(badge.Description() + "\n")
		} else {

			b.WriteString(badge.Description() + "\t")
		}
	}
	//intro
	s_intro := m.Responses[Intro].(string)
	b.WriteString("## Introduction\n")
	b.WriteString(s_intro + "\n")
	//slice of all generic installation options needed for this repo
	xs_install_choices := m.Responses[Install_choices].([]string)
	b.WriteString("## Installation\n")
	//install_users custom install string
	s_install_user := m.Responses[Install_user].(string)
	b.WriteString(s_install_user + "\n")
	//usage
	s_usage := m.Responses[Usage].(string)
	b.WriteString("## Usage\n")
	b.WriteString(s_usage + "\n")
	//contributors

	x_contributors := m.Responses[Contributor_l].([]Contributor)
	b.WriteString("## Contributors\n")
	for _, person := range x_contributors {
		b.WriteString("- " + person.Login + "\n\t- " + person.GitHub + "\n")
	}
	//wild
	s_wild := m.Responses[Wild].(string)
	b.WriteString(s_wild + "\n\n")
	//shameless plug
	s_plug := "This README was created with [Md_Maker](https://github.com/kyleochata/dk_md_make-ah)"
	b.WriteString(s_plug + "\n")
	log.Println(xs_install_choices)
	CreateMD(b.String())
}
func CreateMD(content string) {
	// Create or overwrite the file "test.md"
	file, err := os.Create("test.md")
	if err != nil {
		log.Println("Error creating md file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Println("Error writing to md file:", err)
		return
	}

	log.Println("Markdown file successfully created!")
}
