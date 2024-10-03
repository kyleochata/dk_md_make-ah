package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	lli "github.com/kyleochata/md_maker/licenseitem"
)

type License_model struct {
	Answers
	FoundLicense   bool
	License_type   string
	Pregen_license string
	TextArea       textarea.Model
	edit_pregen    bool
	List           list.Model
	showList       bool
}

func (m License_model) Init() tea.Cmd { return nil }
func (m License_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+b":
			return m.send_to_usage()
		case "tab":
			if m.FoundLicense {
				return m.edit_pregen_license()
			}
		case "ctrl+l":
			return m.show_license_list()
		}
	}
	var cmd tea.Cmd
	if m.edit_pregen {
		if m.TextArea, cmd = m.TextArea.Update(msg); cmd != nil {
			return m, cmd
		}
	}
	if m.showList {
		if m.List, cmd = m.List.Update(msg); cmd != nil {
			return m, cmd
		}
	}

	return m, nil
}
func (m License_model) View() string {

	uiEl := []string{}
	if m.FoundLicense {
		s := fmt.Sprintf("Seems like you already have a %s.\nIf you have any changes you wish to make, press tab.\nCtrl+C for the next section", m.License_type)
		uiEl = append(uiEl, gloss.NewStyle().Width(m.Width).Align(gloss.Center).Render(s))
		if m.edit_pregen {
			uiEl = append(uiEl, m.TextArea.View())
		} else {
			uiEl = append(uiEl, gloss.NewStyle().Align(gloss.Left).Render(m.preGenLicense()))
		}
	} else {
		if !m.showList {
			s_no_license := "Seems like there is no LICENSE file in this directory\nIf you have one, please move it to this directory\nCtrl+R: Refresh | Ctrl+L: Generate a license file"
			uiEl = append(uiEl, gloss.NewStyle().Width(m.Width).Align(gloss.Center).Margin(1, 0, 1, 0).Render(s_no_license))
		} else {
			uiEl = append(uiEl, gloss.NewStyle().Width(m.Width).Margin(1, 0, 1, 0).Align(gloss.Center).Render("Please choose a license you would like to use for this project.\nA license of the type you choose will be created for you at the end of this program."))
			uiEl = append(uiEl, m.List.View())
		}
	}
	return gloss.JoinVertical(gloss.Center, uiEl...)
}

func (m *License_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}

func (m License_model) send_to_usage() (tea.Model, tea.Cmd) {
	return New_Use_model(m.Answers), func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: m.Height,
			Width:  m.Width,
		}
	}
}

func New_License_model(a Answers) tea.Model {
	found := licenseFileExists()
	var l_type string
	if found {
		l_type = findLicenseType()
	}
	ta := textarea.New()
	ta.SetWidth(a.Width)
	ta.Blur()
	return License_model{Answers: a, FoundLicense: found, License_type: l_type, TextArea: ta}
}

func licenseFileExists() bool {
	_, err1 := os.Stat("LICENSE")
	_, err2 := os.Stat("license")
	_, err3 := os.Stat("License")
	if os.IsNotExist(err1) && os.IsNotExist(err2) && os.IsNotExist(err3) {
		return false
	}
	return true
}

// const found_s string = "## Licesne\n ```Please see the LICENSE file in this respository```\n Please click on the badge at the top of the README.md for additional information."
const found_s string = `## License
> License
>
> Please review the LICENSE file in this repository
Please click on the badge at the top of the README.md for additional information.`

func (m *License_model) edit_pregen_license() (tea.Model, tea.Cmd) {
	if !m.edit_pregen {
		m.edit_pregen = !m.edit_pregen
		if m.TextArea.Value() == "" {
			m.TextArea.SetValue(found_s)
		}
		m.TextArea.Focus()
		return m, nil
	}
	m.edit_pregen = !m.edit_pregen
	m.TextArea.Blur()
	return m, nil
}

func (m License_model) preGenLicense() string {
	if m.FoundLicense {
		if m.License_type != "" && m.TextArea.Value() == "" {
			return found_s
		}
		if m.License_type != "" {
			return m.TextArea.Value()
		}
	}
	return ""
}
func findLicenseType() string {
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

func (m *License_model) show_license_list() (tea.Model, tea.Cmd) {
	xs_licenses := []string{"MIT", "Apache 2.0", "BSD 3-Clause", "BSD 2-Clause", "GNU GPL v3", "ISC"}
	xli_licenses := []lli.Item{}
	for _, license := range xs_licenses {
		xli_licenses = append(xli_licenses, lli.New_list_item(license))
	}
	listItems := make([]list.Item, len(xli_licenses))
	for i, license := range xli_licenses {
		listItems[i] = license
	}
	m.List = list.New(listItems, lli.CustomDelegate{}, m.Width, m.Height)
	m.List.SetWidth(m.Width)
	m.List.SetHeight(m.Height * 3 / 5)
	m.toggle_show_list()
	m.List.Title = "Available licenses:"
	return m, nil
}
func (m *License_model) toggle_show_list() {
	m.showList = !m.showList
}

//TODO: Send to next module. Save the data before sending

//TODO: Create a custom help and keymap for the license list

//TODO: Create stylings for all text in this model

// func Generate_License_file(license_type string) {

// }
