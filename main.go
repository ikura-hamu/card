package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"go.ikura-hamu.work/card/internal/about"
	"go.ikura-hamu.work/card/internal/repo"
	"go.ikura-hamu.work/card/internal/tabs"
)

func main() {
	model := tabs.NewTabsManager([]string{"About", "Repositories"}, []tea.Model{
		about.NewModel(),
		repo.NewModel(),
	})
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
