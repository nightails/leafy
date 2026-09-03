package main

import (
	"log"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nightails/leafy/app/tui"
)

// Version update with Git Tag
var version = "dev"

func main() {
	v := getVersion()
	m := tui.New(v)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Panicf("Program exited with error: %v", err)
	}
}

func getVersion() string {
	// CI/CD build
	if version != "dev" {
		return version
	}

	// Read build info from Go Build
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}

	// Fallback to "dev"
	return version
}
