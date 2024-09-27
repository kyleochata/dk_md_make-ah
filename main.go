package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	// p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	p := tea.NewProgram(New_Title_model(Answers{}), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error at launch:\t %v", err)
		os.Exit(1)
	}
}
