package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"go.ikura-hamu.work/card/internal/about"
	"go.ikura-hamu.work/card/internal/repo"
	"go.ikura-hamu.work/card/internal/tabs"
)

func main() {
	model, err := tabs.NewTabsManager([]tabs.Tab{
		about.NewModel(),
		repo.NewModel(),
	})
	if err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
