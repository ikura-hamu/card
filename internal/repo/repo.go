package repo

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/google/go-github/v73/github"
	"go.ikura-hamu.work/card/internal/common/merrors"
)

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
		return merrors.New(fmt.Errorf("fetch user: %w", err))
	}

	repos, _, err := client.Repositories.ListByUser(ctx, githubUserName, &github.RepositoryListByUserOptions{
		Type: "owner",
		Sort: "pushed",
		ListOptions: github.ListOptions{
			PerPage: user.GetPublicRepos(),
		},
	})
	if err != nil {
		return merrors.New(fmt.Errorf("fetch repositories: %w", err))
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

type Model struct {
	ready         bool
	reposViewport viewport.Model
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			m.reposViewport = viewport.New(msg.Width-5, msg.Height-5)
			m.ready = true
		} else {
			m.reposViewport.Width = msg.Width - 5
			m.reposViewport.Height = msg.Height - 5
		}
	case merrors.Msg:
		return m, tea.Quit // Handle error appropriately in a real application
	}

	var cmd tea.Cmd
	m.reposViewport, cmd = m.reposViewport.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	if !m.ready {
		return "Loading repositories..."
	}
	return m.reposViewport.View()
}
