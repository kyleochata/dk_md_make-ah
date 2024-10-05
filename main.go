package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("usage: mdmaker go")
		return
	}

	if os.Args[1] == "go" {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
		p := tea.NewProgram(New_Title_model(Answers{}), tea.WithAltScreen())
		// p := tea.NewProgram(New_Title_model(Answers{}))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error at launch:\t %v", err)
			os.Exit(1)
		}
	}
	fmt.Println("usage: mdmaker go")
}

func SendWindowMsg(height, width int) tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: height,
			Width:  width,
		}
	}
}
