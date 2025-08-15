package repo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/v73/github"
)

type reposRateLimitMsg struct {
	resetAt time.Time
}

const githubUserName = "ikura-hamu"

func fetchRepositories() tea.Msg {
	ctx := context.Background()
	client := github.NewClient(http.DefaultClient)

	user, _, err := client.Users.Get(ctx, githubUserName)
	re := &github.RateLimitError{}
	if rateLimit := errors.As(err, &re); rateLimit {
		return reposRateLimitMsg{resetAt: re.Rate.Reset.Time}
	}
	if err != nil {
		return fmt.Errorf("fetch user: %w", err)
	}

	repos, _, err := client.Repositories.ListByUser(ctx, githubUserName, &github.RepositoryListByUserOptions{
		Type: "owner",
		Sort: "pushed",
		ListOptions: github.ListOptions{
			PerPage: user.GetPublicRepos(),
		},
	})
	if rateLimit := errors.As(err, &re); rateLimit {
		return reposRateLimitMsg{resetAt: re.Rate.Reset.Time}
	}
	if err != nil {
		return fmt.Errorf("fetch repositories: %w", err)
	}

	repoList := make([]repo, 0, len(repos))
	for _, r := range repos {
		if r.GetFork() || r.GetArchived() {
			continue
		}

		repoList = append(repoList, repo{
			name:        r.GetName(),
			description: r.GetDescription(),
			starsCount:  r.GetStargazersCount(),
			language:    r.GetLanguage(),
			pushedAt:    r.GetPushedAt().Time,
			topics:      r.Topics,
			url:         r.GetHTMLURL(),
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
