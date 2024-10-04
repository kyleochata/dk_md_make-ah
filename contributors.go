package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Contributor struct {
	Login string `json:"login"`
	Email string `json:"email,omitempty"`
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
func fetchContributorsWithCLI(repoOwner, repoName string) ([]Contributor, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "collaborators", "--jq", ".collaborators[] | {login: .login, email: .email}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch contributors using GitHub CLI: %w", err)
	}

	var contributors []Contributor
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			var contributor Contributor
			err := json.Unmarshal([]byte(line), &contributor)
			if err != nil {
				return nil, fmt.Errorf("failed to parse contributor: %w", err)
			}
			contributors = append(contributors, contributor)
		}
	}
	return contributors, nil
}

// If no Github CLI, http.Get
func fetchContributorsFromAPI(repoOwner, repoName string) ([]Contributor, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contributors", repoOwner, repoName)

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

	var contributors []Contributor
	err = json.Unmarshal(body, &contributors)
	if err != nil {
		return nil, err
	}

	return contributors, nil
}

func fetchContributorsCmd() tea.Cmd {
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
			contributors, err = fetchContributorsWithCLI(repoOwner, repoName)
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
	return fetchContributorsCmd()
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
	}
	return m, nil
}
func (m Contributors_model) View() string {
	return "from contributors"
}
func (m *Contributors_model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.Height, m.Width = msg.Height, msg.Width
}
func New_Contributors_model(a Answers) tea.Model {
	return Contributors_model{Answers: a}
}
