package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	cli "github.com/kyleochata/md_maker/licenseitem"
)

type Has_License_model struct {
	Answers
	licenseType string
	content     string
	TextArea    textarea.Model
	editContent bool //show ta or not
	makeFile    bool
}

func (m Has_License_model) Init() tea.Cmd { return nil }
func (m Has_License_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.toggleTAFocus()
			return m, nil
		}
	}
	if m.editContent {
		var cmd tea.Cmd
		if m.TextArea, cmd = m.TextArea.Update(msg); cmd != nil {
			return m, cmd
		}
	}
	return m, nil
}
func (m Has_License_model) View() string {
	var b strings.Builder
	b.WriteString(m.content + "\n")
	b.WriteString(m.licenseType)

	return b.String()
}

func (m *Has_License_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	m.TextArea.SetWidth(m.Width)
}

func (m *Has_License_model) toggleTAFocus() {
	m.editContent = !m.editContent
	if !m.editContent {
		m.TextArea.Blur()
	} else {
		m.TextArea.Focus()
	}
}

func New_has_License_model(a Answers, l_type string, make bool) tea.Model {
	ta := textarea.New()
	ta.SetWidth(a.Width)
	ta.Focus()
	content := readmeLicenseContent(l_type)
	prevContent, ok := a.Responses["licenseContent"].(string)
	if ok && prevContent != "" {
		ta.SetValue(prevContent)
	} else {
		ta.SetValue(content)
	}

	return Has_License_model{Answers: a, TextArea: ta, editContent: false, licenseType: l_type, makeFile: make, content: content}
}
func readmeLicenseContent(lt string) string {
	var content string = ""
	if lt != "" {
		content = fmt.Sprintf("##%s\n>\n> Please review the %s file in this repository\nPlease click on the badge at the top of the README for additional information.", lt, strings.ToUpper(lt))
	}
	return content
}

type Fail_License_check_model struct {
	Answers
}

func (m Fail_License_check_model) Init() tea.Cmd { return nil }
func (m Fail_License_check_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			return m.retryLicenseCheck()
		case "tab":
			return m.send_to_Available_license_model()
		}
	}
	return m, nil
}
func (m Fail_License_check_model) View() string {
	var sb strings.Builder
	sb.WriteString("No License/license/LICENSE file detected\n")
	sb.WriteString("If you have a file with a license, please move it to the root of the current working directory.\n")
	sb.WriteString("Please ensure that it is named: License, LICENSE, or license.\n")
	sb.WriteString("Press Tab to create a LICENSE\n")
	sb.WriteString("Ctrl+C: Quit | Ctrl+N: Advance to next section | Tab: Create License")
	return sb.String()
}
func (m *Fail_License_check_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}

func New_Fail_License_check_model(a Answers) tea.Model {
	return Fail_License_check_model{Answers: a}
}
func (m Fail_License_check_model) retryLicenseCheck() (tea.Model, tea.Cmd) {
	if LicenseFileExists() {
		return New_has_License_model(m.Answers, FindLicenseType(), false), func() tea.Msg {
			return tea.WindowSizeMsg{
				Height: m.Height,
				Width:  m.Width,
			}
		}
	}
	return m, nil
}

func (m Fail_License_check_model) send_to_Available_license_model() (tea.Model, tea.Cmd) {
	return new_available_license_model(m.Answers), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: m.Height,
			Width:  m.Width,
		}
	}
}

type available_license_model struct {
	Answers
	List list.Model
}

func (m available_license_model) Init() tea.Cmd { return nil }
func (m available_license_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+b":
			return New_Fail_License_check_model(m.Answers), nil
		case "enter":
			return New_license_from_list_model(m.Answers, m.List.SelectedItem().FilterValue())
		}
	}
	var cmd tea.Cmd
	if m.List, cmd = m.List.Update(msg); cmd != nil {
		return m, cmd
	}
	return m, nil
}
func (m available_license_model) View() string {
	return m.List.View()
}
func (m *available_license_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}
func new_available_license_model(a Answers) tea.Model {
	litems := loadLicenses()
	log.Println(litems)
	maxListHeight := a.Height / 3
	if maxListHeight < 12 {
		maxListHeight = 12
	}
	l := list.New(litems, cli.CustomDelegate{}, a.Width-2, maxListHeight)
	l.Title = "Available Licenses"
	return available_license_model{Answers: a, List: l}
}

func New_license_from_list_model(a Answers, l_type string) (tea.Model, tea.Cmd) {
	return New_has_License_model(a, l_type, true), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: a.Height,
			Width:  a.Width,
		}
	}
}

func loadLicenses() []list.Item {
	licenses, _ := cli.GetAvailableLicenses()
	listItems := make([]list.Item, len(licenses))
	for i, license := range licenses {
		listItems[i] = license
	}
	return listItems
}

func LicenseFileExists() bool {
	_, err1 := os.Stat("LICENSE")
	_, err2 := os.Stat("license")
	_, err3 := os.Stat("License")
	if os.IsNotExist(err1) && os.IsNotExist(err2) && os.IsNotExist(err3) {
		return false
	}
	return true
}

func FindLicenseType() string {
	// Try to open the file with possible variations of the license filename
	var file *os.File
	var err error

	// Check for "LICENSE", "license", or "License"
	for _, name := range []string{"LICENSE", "license", "License"} {
		file, err = os.Open(name)
		if err == nil {
			break // If the file opens successfully, exit the loop
		}
	}

	if err != nil {
		return "No license file found"
	}
	defer file.Close()

	// Read the file content to find the license type
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text()) // Convert to lowercase for case-insensitive matching
		if strings.Contains(line, "mit license") {
			return "MIT License"
		}
		if strings.Contains(line, "gnu general public license") || strings.Contains(line, "gpl") {
			return "GNU General Public License"
		}
		if strings.Contains(line, "apache license") {
			return "Apache License"
		}
	}

	if err := scanner.Err(); err != nil {
		return "Error reading license file"
	}
	return ""
}
