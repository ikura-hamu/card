package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"go.ikura-hamu.work/card/internal/about"
	"go.ikura-hamu.work/card/internal/icon"
	"go.ikura-hamu.work/card/internal/repo"
	"go.ikura-hamu.work/card/internal/tabs"
)

func run() error {
	defer icon.CleanupTempIcons()

	model, err := tabs.NewTabsManager([]tabs.Tab{
		about.NewModel(),
		repo.NewModel(),
		icon.NewModel(),
	})
	if err != nil {
		return fmt.Errorf("failed to create tabs manager: %w", err)
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
