package about

import (
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"go.ikura-hamu.work/card/internal/common/merrors"
)

type readMeMsg struct {
	content string
}

func fetchReadme() (msg tea.Msg) {
	resp, err := http.Get("https://raw.githubusercontent.com/ikura-hamu/ikura-hamu/refs/heads/main/README.md")
	if err != nil {
		return fmt.Errorf("fetch readme: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			msg = fmt.Errorf("close resp.Body: %w", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetch readme: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read readme: %w", err)
	}

	return readMeMsg{content: string(body)}
}

type Model struct {
	ready      bool
	mdViewport viewport.Model
}

func NewModel() Model {
	return Model{
		ready:      false,
		mdViewport: viewport.New(0, 0),
	}
}

func (m Model) Name() string {
	return "About"
}

func (m Model) Init() tea.Cmd {
	return fetchReadme
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case readMeMsg:
		md, err := glamour.Render(msg.content, "dark")
		if err != nil {
			return m, merrors.NewCmd(fmt.Errorf("render markdown: %w", err))
		}
		m.mdViewport.SetContent(md)
	case tea.WindowSizeMsg:
		if !m.ready {
			m.mdViewport = viewport.New(msg.Width, msg.Height)
			m.ready = true
		} else {
			m.mdViewport.Width = msg.Width
			m.mdViewport.Height = msg.Height
		}
	}

	var cmd tea.Cmd
	m.mdViewport, cmd = m.mdViewport.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "Loading readme..."
	}
	return m.mdViewport.View()
}

func (m Model) KeyMap() help.KeyMap {
	return KeyMap(m.mdViewport.KeyMap)
}
