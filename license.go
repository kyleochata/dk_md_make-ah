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
	gloss "github.com/charmbracelet/lipgloss"
	cli "github.com/kyleochata/md_maker/licenseitem"
)

const (
	License_t string = "licenseType"
	License_c string = "licenseContent"
	License_w string = "createLicenseFile"
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
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.toggleTAFocus()
			return m, nil
		case "ctrl+l":
			m.Responses["licenseType"] = m.licenseType
			return new_available_license_model(m.Answers), nil
		case "ctrl+n":
			if !m.editContent {
				m.save_license_data()
				return Send_to_Contributors(m.Answers)
			}
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
	b.WriteString(gloss.NewStyle().Margin(1, 0, 1, 0).Render("Edit License section of README\n"))
	b.WriteString(gloss.NewStyle().Margin(1, 0, 1, 0).Render(m.TextArea.View()))
	b.WriteString("Tab: Review Changes | Ctrl+L: Change License |Ctrl+C: Quit")
	if !m.editContent {
		b.Reset()
		b.WriteString(fmt.Sprintf("Found a %s type in your current working directory\n", m.licenseType))
		b.WriteString("The below content will be added to the License section of the README\n")
		b.WriteString(m.contentStyle() + "\n")
		b.WriteString("Tab: Edit license content | Ctrl+L: Change license type\n")
		b.WriteString("Ctrl+C: Quit | Ctrl+N: Next section")
	}
	return gloss.NewStyle().Width(m.Width).Align(gloss.Center).Render(b.String())
}
func (m *Has_License_model) save_license_data() {
	m.Responses[License_t] = m.licenseType
	m.Responses[License_c] = m.content
	m.Responses[License_w] = m.makeFile
}
func Send_to_Contributors(a Answers) (tea.Model, tea.Cmd) {
	return New_Contributors_model(a), SendWindowMsg(a.Height, a.Width)
}

func (m *Has_License_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	m.TextArea.SetWidth(m.Width)
}

func (m *Has_License_model) toggleTAFocus() {
	if m.editContent {
		m.content = m.TextArea.Value()
		m.TextArea.Blur()
	} else {
		m.TextArea.Focus()
	}
	m.editContent = !m.editContent
}

func New_has_License_model(a Answers, l_type string, make bool) tea.Model {
	ta := textarea.New()
	ta.SetWidth(a.Width)
	ta.Focus()
	prev_l_type, ok := a.Responses[License_t].(string)
	if ok && prev_l_type != "" {
		l_type = prev_l_type
	}
	prevContent, ok := a.Responses[License_c].(string)
	var content string
	if ok && prevContent != "" {
		ta.SetValue(prevContent)
	} else {
		content = readmeLicenseContent(l_type)
		ta.SetValue(content)
	}
	return Has_License_model{Answers: a, TextArea: ta, editContent: false, licenseType: l_type, makeFile: make, content: ta.Value()}
}
func readmeLicenseContent(lt string) string {
	var content string = ""
	if lt != "" {
		content = fmt.Sprintf("##License\n\n> %s\n>\n> Please review the %s file in this repository\nPlease click on the badge at the top of the README for additional information.", lt, strings.ToUpper(lt))
	}
	return content
}

func (m Has_License_model) contentStyle() string {
	return gloss.NewStyle().Align(gloss.Left).Border(gloss.RoundedBorder()).Width(m.Width-4).Margin(1, 0, 1, 0).Render(m.content)
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
		return m, tea.ClearScreen
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
	return gloss.NewStyle().Height(m.Height).Render(m.List.View())
}
func (m *available_license_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	maxListHeight := m.Height / 3
	if maxListHeight < 12 {
		maxListHeight = 12
	}
}
func new_available_license_model(a Answers) tea.Model {
	litems := LoadLicenses()
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
	writeNewLicense := true
	prev_l_type := a.Responses[License_t].(string)
	if prev_l_type == l_type {
		writeNewLicense = false
	}
	log.Printf("Chosen license type: %s", l_type)
	log.Printf("Previous license type: %s", prev_l_type)
	log.Printf("WriteNewLicense: %v", writeNewLicense)
	return New_has_License_model(a, l_type, writeNewLicense), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: a.Height,
			Width:  a.Width,
		}
	}
}

func LoadLicenses() []list.Item {
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
