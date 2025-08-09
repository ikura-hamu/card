package repo

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/cli/browser"
	"go.ikura-hamu.work/card/internal/common/merrors"
	"go.ikura-hamu.work/card/internal/common/size"
)

type repo struct {
	name        string
	description string
	starsCount  int
	language    string
	pushedAt    time.Time
	topics      []string
	url         string
}

var _ list.Item = repo{}

func (r repo) Title() string       { return r.name }
func (r repo) Description() string { return r.description }
func (r repo) FilterValue() string { return fmt.Sprintf("%s %s %s", r.name, r.description, r.language) }

type reposMsg struct {
	repos []repo
}

type Model struct {
	ready         bool
	reposViewport viewport.Model
	repoList      list.Model
	size          size.Size
}

func NewModel() Model {
	return Model{
		ready:         false,
		reposViewport: viewport.New(0, 0),
	}
}

func (m Model) Name() string {
	return "Repositories"
}

func (m Model) Init() tea.Cmd {
	return fetchRepositories
}

func additionalListKeyBindings() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys("o", tea.KeyEnter.String()),
			key.WithHelp("o", "open in browser"),
		),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 3)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String(), "o":
			currentRepo := m.repoList.SelectedItem().(repo)
			err := browser.OpenURL(currentRepo.url)
			if err != nil {
				return m, merrors.NewCmd(fmt.Errorf("open repository URL: %w", err))
			}
		}
	case reposMsg:
		items := make([]list.Item, 0, len(msg.repos))
		for _, r := range msg.repos {
			items = append(items, r)
		}
		cmd := m.repoList.SetItems(items)
		cmds = append(cmds, cmd)
	case reposRateLimitMsg:
		content := fmt.Sprintf("Rate limit exceeded. Please wait until %s to try again.", msg.resetAt.Format(time.DateTime))
		md, err := glamour.Render(content, "dark")
		if err != nil {
			return m, tea.Quit
		}
		m.reposViewport.SetContent(md)
	case tea.WindowSizeMsg:
		m.size = size.Size{Width: msg.Width, Height: msg.Height}
		listSize := repoListSize(m.size)
		mdTextViewSize := repoMDTextViewSize(m.size)
		if !m.ready {
			m.reposViewport = viewport.New(mdTextViewSize.Width, mdTextViewSize.Height)

			m.repoList = list.New([]list.Item{}, list.NewDefaultDelegate(), listSize.Width, listSize.Height)
			m.repoList.AdditionalFullHelpKeys = additionalListKeyBindings
			m.repoList.AdditionalShortHelpKeys = additionalListKeyBindings
			m.repoList.Title = "Repositories"

			m.ready = true
		} else {
			m.reposViewport.Width = mdTextViewSize.Width
			m.reposViewport.Height = mdTextViewSize.Height
			m.repoList.SetSize(listSize.Width, listSize.Height)
		}
	case merrors.Msg:
		return m, tea.Quit // Handle error appropriately in a real application
	}

	var cmd tea.Cmd
	m.reposViewport, cmd = m.reposViewport.Update(msg)
	cmds = append(cmds, cmd)
	m.repoList, cmd = m.repoList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func repoListSize(s size.Size) size.Size {
	return size.Size{
		Width:  s.Width / 2,
		Height: s.Height,
	}
}

func repoMDTextViewSize(s size.Size) size.Size {
	x, y := mdContainerStyle.GetFrameSize()
	return size.Size{
		Width:  s.Width/2 - x,
		Height: s.Height - y,
	}
}

func repoMarkdownView(r repo) (string, error) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## %s\n\n", r.Title()))
	b.WriteString(fmt.Sprintf("%s\n\n", r.Description()))
	b.WriteString(fmt.Sprintf("- Stars: %d\n", r.starsCount))
	b.WriteString(fmt.Sprintf("- Language: %s\n", r.language))
	b.WriteString(fmt.Sprintf("- Last pushed: %s\n", r.pushedAt.Format(time.RFC1123)))
	b.WriteString(fmt.Sprintf("- Topics: %s\n", strings.Join(r.topics, ", ")))

	md, err := glamour.Render(b.String(), "dark")
	if err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	return md, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Loading repositories..."
	}

	currentRepo, ok := m.repoList.SelectedItem().(repo)
	if !ok {
		return "No repository selected."
	}

	repoMD, err := repoMarkdownView(currentRepo)
	if err != nil {
		return fmt.Sprintf("Error rendering repository details: %v", err)
	}
	m.reposViewport.SetContent(repoMD)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		listContainerStyle(repoListSize(m.size)).Render(m.repoList.View()),
		mdContainerStyle.Render(m.reposViewport.View()))
}

func listContainerStyle(cs size.Size) lipgloss.Style {
	return listContainerBaseStyle.Width(cs.Width).Height(cs.Height)
}

var (
	listContainerBaseStyle = lipgloss.NewStyle()
	mdContainerStyle       = lipgloss.NewStyle().
				Border(lipgloss.ASCIIBorder())
)
