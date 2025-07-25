package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/google/go-github/v73/github"
)

type model struct {
	readMeContent string
	mdViewport    viewport.Model
	ready         bool
	reposViewport viewport.Model
	tabHeader     string

	tabs   []string
	active int
}

type errMsg struct {
	err error
}

type readMeMsg struct {
	content string
}

type repo struct {
	name        string
	desctiption string
	starsCount  int
	language    string
	pushedAt    time.Time
}

type reposMsg struct {
	repos []repo
}

const githubUserName = "ikura-hamu"

func fetchRepositories() tea.Msg {
	ctx := context.Background()
	client := github.NewClient(http.DefaultClient)

	user, _, err := client.Users.Get(ctx, githubUserName)
	if err != nil {
		return errMsg{err: fmt.Errorf("fetch user: %w", err)}
	}

	repos, _, err := client.Repositories.ListByUser(ctx, githubUserName, &github.RepositoryListByUserOptions{
		Type: "owner",
		Sort: "pushed",
		ListOptions: github.ListOptions{
			PerPage: user.GetPublicRepos(),
		},
	})
	if err != nil {
		return errMsg{err: fmt.Errorf("fetch repositories: %w", err)}
	}

	repoList := make([]repo, 0, len(repos))
	for _, r := range repos {
		if r.GetFork() || r.GetArchived() {
			continue
		}

		repoList = append(repoList, repo{
			name:        r.GetName(),
			desctiption: r.GetDescription(),
			starsCount:  r.GetStargazersCount(),
			language:    r.GetLanguage(),
			pushedAt:    r.GetPushedAt().Time,
		})
	}

	slices.SortFunc(repoList, func(a, b repo) int {
		if a.starsCount != b.starsCount {
			return b.starsCount - a.starsCount // Sort by stars count descending
		}

		return b.pushedAt.Compare(a.pushedAt) // Sort by pushed date descending
	})

	return reposMsg{
		repos: repoList,
	}
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
	return tea.Batch(fetchReadme, fetchRepositories)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.active++
			m.active %= len(m.tabs)
			m.tabHeader = m.tabs[m.active]
		case "shift+tab":
			m.active--
			if m.active < 0 {
				m.active += len(m.tabs)
			}
			m.active %= len(m.tabs)
			m.tabHeader = m.tabs[m.active]
		}
	case errMsg:
		return m, tea.Quit
	case readMeMsg:
		m.readMeContent = msg.content
		md, err := glamour.Render(m.readMeContent, "dark")
		if err != nil {
			return m, tea.Quit
		}
		m.mdViewport.SetContent(md)
	case reposMsg:
		contentBdr := strings.Builder{}
		contentBdr.WriteString("# ikura-hamu Repositories\n\n")
		for _, r := range msg.repos {
			contentBdr.WriteString(fmt.Sprintf(`## %s
Description: %s	

⭐: %d	

Language: %s

`,
				r.name, r.desctiption, r.starsCount, r.language))
		}
		md, err := glamour.Render(contentBdr.String(), "dark")
		if err != nil {
			return m, tea.Quit
		}
		m.reposViewport.SetContent(md)
	case tea.WindowSizeMsg:
		if !m.ready {
			m.mdViewport = viewport.New(msg.Width, msg.Height-2)
			m.mdViewport.SetContent(m.readMeContent)
			m.reposViewport = viewport.New(msg.Width, msg.Height-2)
			m.reposViewport.SetContent(m.reposViewport.View())
			m.ready = true
		} else {
			m.mdViewport.Width = msg.Width
			m.mdViewport.Height = msg.Height
		}
	}

	var mdCmd, repoCmd tea.Cmd
	m.mdViewport, mdCmd = m.mdViewport.Update(msg)
	m.reposViewport, repoCmd = m.reposViewport.Update(msg)

	return m, tea.Batch(mdCmd, repoCmd)
}

func (m model) View() string {
	header := fmt.Sprintf("Tab: %s\n", m.tabHeader)
	if m.tabs[m.active] == "about" {
		return header + m.mdViewport.View()
	}
	if m.tabs[m.active] == "repos" {
		return header + m.reposViewport.View()
	}

	return "unknown tab\n"
}

func main() {
	m := model{
		tabs:      []string{"about", "repos"},
		tabHeader: "about",
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
