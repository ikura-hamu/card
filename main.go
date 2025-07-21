package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type model struct {
	readMeContent string
	viewport      viewport.Model
	ready         bool
}

type errMsg struct {
	err error
}

type readMeMsg struct {
	content string
}

func fetchReadme() tea.Msg {
	resp, err := http.Get("https://raw.githubusercontent.com/ikura-hamu/ikura-hamu/refs/heads/main/README.md")
	if err != nil {
		return errMsg{err: fmt.Errorf("fetch readme: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errMsg{err: fmt.Errorf("fetch readme: %w", err)}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errMsg{err: fmt.Errorf("read readme: %w", err)}
	}

	return readMeMsg{content: string(body)}
}

func (m model) Init() tea.Cmd {
	return fetchReadme
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case errMsg:
		return m, tea.Quit
	case readMeMsg:
		m.readMeContent = msg.content
		md, err := glamour.Render(m.readMeContent, "dark")
		if err != nil {
			return m, tea.Quit
		}
		m.viewport.SetContent(md)
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.SetContent(m.readMeContent)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.readMeContent == "" {
		return "Loading README...\n"
	}

	return m.viewport.View()
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
