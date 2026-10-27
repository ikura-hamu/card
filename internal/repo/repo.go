package repo

import (
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/cli/browser"
	"go.ikura-hamu.work/card/internal/common/merrors"
	"go.ikura-hamu.work/card/internal/common/size"
	"go.ikura-hamu.work/card/internal/common/theme"
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
	repoList := list.New([]list.Item{}, themedDelegate(), 0, 0)
	repoList.AdditionalFullHelpKeys = additionalListKeyBindings
	repoList.AdditionalShortHelpKeys = additionalListKeyBindings
	repoList.Title = "Repositories"
	styleRepositoryList(&repoList)
	repoList.SetShowHelp(false)

	reposViewport := viewport.New()
	reposViewport.SoftWrap = true

	return Model{
		ready:         false,
		reposViewport: reposViewport,
		repoList:      repoList,
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
			key.WithKeys("o", "enter"),
			key.WithHelp("o", "open in browser"),
		),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 3)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter", "o":
			currentRepo, ok := m.repoList.SelectedItem().(repo)
			if !ok {
				break
			}
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
		md, err := theme.RenderMarkdown(content, 80)
		if err != nil {
			return m, merrors.NewCmd(fmt.Errorf("render markdown: %w", err))
		}
		m.reposViewport.SetContent(md)
	case tea.WindowSizeMsg:
		m.size = size.Size{Width: msg.Width, Height: msg.Height}
		listSize := repoListSize(m.size)
		mdTextViewSize := repoMDTextViewSize(m.size)
		m.reposViewport.SetWidth(mdTextViewSize.Width)
		m.reposViewport.SetHeight(mdTextViewSize.Height)
		m.repoList.SetSize(listSize.Width, listSize.Height)
		m.ready = true
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
	listWidth := repoListSize(s).Width
	return size.Size{
		Width:  max(0, s.Width-listWidth-x),
		Height: max(0, s.Height-y),
	}
}

func repoMarkdownView(r repo, width int) (string, error) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## %s\n\n", r.Title()))
	b.WriteString(fmt.Sprintf("%s\n\n", r.Description()))
	b.WriteString(fmt.Sprintf("- Stars: %d\n", r.starsCount))
	b.WriteString(fmt.Sprintf("- Language: %s\n", r.language))
	b.WriteString(fmt.Sprintf("- Last pushed: %s\n", r.pushedAt.Format(time.RFC1123)))
	b.WriteString(fmt.Sprintf("- Topics: %s\n", strings.Join(r.topics, ", ")))

	md, err := theme.RenderMarkdown(b.String(), width)
	if err != nil {
		return "", fmt.Errorf("render markdown: %w", err)
	}
	return md, nil
}

func (m Model) View() tea.View {
	if !m.ready {
		return tea.NewView("Loading repositories...")
	}

	currentRepo, ok := m.repoList.SelectedItem().(repo)
	if !ok {
		return tea.NewView("No repository selected.")
	}

	repoMD, err := repoMarkdownView(currentRepo, m.reposViewport.Width())
	if err != nil {
		return tea.NewView(fmt.Sprintf("Error rendering repository details: %v", err))
	}
	m.reposViewport.SetContent(repoMD)

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top,
		listContainerStyle(repoListSize(m.size)).Render(m.repoList.View()),
		mdContainerStyle.Render(m.reposViewport.View())))
}

func (m Model) KeyMap() help.KeyMap {
	return m.repoList
}

func (m Model) Help() help.Model {
	return m.repoList.Help
}

func listContainerStyle(cs size.Size) lipgloss.Style {
	return listContainerBaseStyle.
		Width(cs.Width).
		Height(cs.Height).
		MaxWidth(cs.Width).
		MaxHeight(cs.Height)
}

var (
	listContainerBaseStyle = lipgloss.NewStyle()
	mdContainerStyle       = lipgloss.NewStyle().
				Border(lipgloss.ASCIIBorder()).
				BorderForeground(lipgloss.Color(theme.Secondary))
)
