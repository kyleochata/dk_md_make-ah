package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Contributor struct {
	Login  string `json:"login"`
	GitHub string `json:"html_url"`
}

// type contributorsMsg struct {
// 	contributors []Contributor
// 	owner        string
// }

type Contributors_model struct {
	Answers
	contributors []Contributor
	// errorMessage string
	table table.Model
	owner string
}

// type errorMsg struct {
// 	err error
// }

func isOnline() bool {
	_, err := net.DialTimeout("tcp", "github.com:80", 3*time.Second)
	return err == nil
}

func isGHCLIInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// Get the repository owner and name from the local git config
func getRepoOwnerAndName() (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Failed to get github repo url")
		return "", "", fmt.Errorf("failed to get GitHub repo URL: %w", err)
	}
	// Parse the URL to extract the owner and repo name
	url := strings.TrimSpace(string(output))
	re := regexp.MustCompile(`github\.com[:/](.+)/(.+?)(\.git)?$`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("could not parse GitHub URL: %s", url)
	}

	return matches[1], matches[2], nil
}

// Get contributors using Github CLI if installed
func fetchContributorsWithCLI() ([]Contributor, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "collaborators", "--jq", ".collaborators[] | {login: .login, email: .email}")
	output, err := cmd.Output()
	if err != nil {
		log.Println("failed to fetch with github cli")
		return nil, fmt.Errorf("failed to fetch contributors using GitHub CLI: %w", err)
	}

	var contributors []Contributor
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			var contributor Contributor
			err := json.Unmarshal([]byte(line), &contributor)
			if err != nil {
				log.Print("failed to parse contributors from gh cli")
				return nil, fmt.Errorf("failed to parse contributor: %w", err)
			}
			contributors = append(contributors, contributor)
		}
	}
	return contributors, nil
}

// If no Github CLI, http.Get
func fetchContributorsFromAPI(repoOwner, repoName string) ([]Contributor, error) {
	// url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", repoOwner, repoName)
	url := "https://api.github.com/repos/kyleochata/Will-DO-Crush-your-goals/contributors"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// log.Println("resp body:", body)
	var contributors []Contributor
	err = json.Unmarshal(body, &contributors)
	// log.Println("Contributors:", contributors)

	if err != nil {
		return nil, err
	}

	return contributors, nil
}

func (m *Contributors_model) FetchContributorsCmd() {
	if !isOnline() {
		log.Println("Not connected to internet")
		return
	}

	repoOwner, repoName, err := getRepoOwnerAndName()
	if err != nil {
		log.Println("Unable to get repoOwner and name")
		return
	}

	var contributors []Contributor
	if isGHCLIInstalled() {
		contributors, err = fetchContributorsWithCLI()
	} else {
		contributors, err = fetchContributorsFromAPI(repoOwner, repoName)
	}

	if err != nil {
		log.Println("Unable to use API or CLI to get Github contributor info")
		return
	}
	m.contributors = contributors
	m.owner = repoOwner

}

func (m Contributors_model) Init() tea.Cmd {
	return nil
}
func (m Contributors_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
		// default:
		// 	log.Println("Msg: ", msg)
	}
	return m, nil
}
func (m Contributors_model) View() string {
	uiEl := []string{gloss.NewStyle().Width(m.Width).Render("from contributors")}

	uiEl = append(uiEl, gloss.NewStyle().Width(m.Width).Render(m.table.View()))
	return gloss.JoinVertical(gloss.Center, uiEl...)
}
func (m *Contributors_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
	m.table.SetHeight(m.Height)
	m.table.SetWidth(m.Width - 4)
}
func New_Contributors_model(a Answers) tea.Model {
	t := table.New()
	xs_column_header := []string{"Index", "Username", "Owner", "GitHub URL"}
	x_tableCol_header := make([]table.Column, len(xs_column_header))
	for i, header := range xs_column_header {
		x_tableCol_header[i] = table.Column{Title: header, Width: (a.Width / len(xs_column_header)) - 4}
	}

	t.SetColumns(x_tableCol_header)
	t.SetHeight(a.Height / 3)
	t.SetWidth(a.Width - 2)
	model := Contributors_model{Answers: a, table: t}
	model.FetchContributorsCmd()
	model.popTableRows()
	log.Println("Owner: ", model.owner)
	return model
}
func (m *Contributors_model) popTableRows() {
	var rows []table.Row
	for i, contributor := range m.contributors {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			contributor.Login,
			strconv.FormatBool(contributor.Login == m.owner),
			contributor.GitHub,
		})
	}
	m.table.SetRows(rows)
}

// func (m Contributors_model) contributorsEqual(newContributors []Contributor) bool {
// 	// Check if the lengths of the slices are different
// 	if len(m.contributors) != len(newContributors) {
// 		return false
// 	}

// 	// Compare contents of the slices
// 	for i := range m.contributors {
// 		if m.contributors[i] != newContributors[i] {
// 			return false
// 		}
// 	}
// 	return true
// }
