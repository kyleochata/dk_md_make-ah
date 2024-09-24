package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	p := tea.NewProgram(New_Title_model(Answers{}), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error at launch:\t %v", err)
		os.Exit(1)
	}
}
