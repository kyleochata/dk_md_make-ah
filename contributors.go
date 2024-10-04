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
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type Contributor struct {
	Login  string `json:"login"`
	GitHub string `json:"html_url"`
}
type contributorsMsg struct {
	contributors []Contributor
}

type Contributors_model struct {
	Answers
	contributors []Contributor
	errorMessage string
}
type errorMsg struct {
	err error
}

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
	log.Println("resp body:", body)
	var contributors []Contributor
	err = json.Unmarshal(body, &contributors)
	log.Println("Contributors:", contributors)

	if err != nil {
		return nil, err
	}

	return contributors, nil
}

func FetchContributorsCmd() tea.Cmd {
	return func() tea.Msg {
		if !isOnline() {
			return errorMsg{err: fmt.Errorf("no internet connection")}
		}

		repoOwner, repoName, err := getRepoOwnerAndName()
		if err != nil {
			return errorMsg{err: err}
		}

		var contributors []Contributor
		if isGHCLIInstalled() {
			contributors, err = fetchContributorsWithCLI()
		} else {
			contributors, err = fetchContributorsFromAPI(repoOwner, repoName)
		}

		if err != nil {
			return errorMsg{err: err}
		}

		return contributorsMsg{contributors: contributors}
	}
}

func (m Contributors_model) Init() tea.Cmd {
	return nil
}
func (m Contributors_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case contributorsMsg:
		m.contributors = msg.contributors
		return m, tea.ClearScreen
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
func (m Contributors_model) View() string {
	uiEl := []string{gloss.NewStyle().Width(m.Width).Render("from contributors")}
	for _, contributor := range m.contributors {
		uiEl = append(uiEl, gloss.NewStyle().Align(gloss.Left).Render(gloss.JoinHorizontal(gloss.Left, contributor.Login+"\t", contributor.GitHub)))
	}
	return gloss.JoinVertical(gloss.Center, uiEl...)
}
func (m *Contributors_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}
func New_Contributors_model(a Answers) tea.Model {
	return Contributors_model{Answers: a}
}
