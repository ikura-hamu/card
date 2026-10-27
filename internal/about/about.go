package about

import (
	"fmt"
	"io"
	"net/http"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"go.ikura-hamu.work/card/internal/common/merrors"
	"go.ikura-hamu.work/card/internal/common/theme"
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
		mdViewport: viewport.New(),
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
		md, err := theme.RenderMarkdown(msg.content, 80)
		if err != nil {
			return m, merrors.NewCmd(fmt.Errorf("render markdown: %w", err))
		}
		m.mdViewport.SetContent(md)
	case tea.WindowSizeMsg:
		if !m.ready {
			m.mdViewport = viewport.New(
				viewport.WithWidth(msg.Width),
				viewport.WithHeight(msg.Height),
			)
			m.ready = true
		} else {
			m.mdViewport.SetWidth(msg.Width)
			m.mdViewport.SetHeight(msg.Height)
		}
	}

	var cmd tea.Cmd
	m.mdViewport, cmd = m.mdViewport.Update(msg)

	return m, cmd
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading readme...")
	}
	return tea.NewView(m.mdViewport.View())
}

func (m Model) KeyMap() help.KeyMap {
	return KeyMap(m.mdViewport.KeyMap)
}
