package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/app/tui"
)

// version of the app. Manually updated.
const version = "v0.2.2"

func main() {
	m := tui.New(version)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}
