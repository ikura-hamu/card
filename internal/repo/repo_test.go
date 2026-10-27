package repo

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"go.ikura-hamu.work/card/internal/tabs"
)

func TestModelUpdateBeforeWindowSize(t *testing.T) {
	tests := []struct {
		name string
		msg  tea.Msg
	}{
		{name: "terminal capability message", msg: struct{}{}},
		{name: "enter without selection", msg: tea.KeyPressMsg{Code: tea.KeyEnter}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.Update(tt.msg)
		})
	}
}

func TestViewUsesAvailableSize(t *testing.T) {
	const height = 19

	for _, width := range []int{40, 41, 78, 79, 120, 160} {
		t.Run(fmt.Sprintf("width %d", width), func(t *testing.T) {
			model := NewModel()
			updated, _ := model.Update(reposMsg{repos: []repo{
				{
					name:        "summer-in-japan-is-too-hot-mcp",
					description: "日本の夏は暑すぎる！",
					language:    "Go",
					pushedAt:    time.Date(2026, time.April, 30, 13, 0, 0, 0, time.UTC),
					topics:      []string{"compiler", "go", "golang"},
				},
				{name: "q-cli", description: "traQ webhookを使ったメッセージ送信CLI"},
				{name: "c3", description: "\"Cleanup Context Checker\" Go static analyzer"},
			}})
			model = updated.(Model)
			updated, _ = model.Update(tea.WindowSizeMsg{Width: width, Height: height})
			model = updated.(Model)

			view := model.View().Content
			if got := lipgloss.Width(view); got != width {
				t.Errorf("view width = %d, want %d", got, width)
			}
			if got := lipgloss.Height(view); got != height {
				t.Errorf("view height = %d, want %d", got, height)
			}
		})
	}
}

func TestDetailsWrapToPaneWidth(t *testing.T) {
	model := NewModel()
	updated, _ := model.Update(reposMsg{repos: []repo{{
		name:        "card",
		description: "A terminal application with a long description that should wrap inside the repository details pane.",
		pushedAt:    time.Date(2026, time.April, 30, 13, 0, 0, 0, time.UTC),
		topics:      []string{"compiler", "go", "golang"},
	}}})
	model = updated.(Model)

	for _, width := range []int{78, 120, 79} {
		updated, _ = model.Update(tea.WindowSizeMsg{Width: width, Height: 40})
		model = updated.(Model)
		var details strings.Builder
		for _, line := range strings.Split(model.View().Content, "\n") {
			line = ansi.Strip(ansi.Cut(line, width/2+1, width-1))
			details.WriteString(line)
			details.WriteByte(' ')
		}
		text := strings.Join(strings.Fields(details.String()), " ")
		for _, want := range []string{
			"A terminal application with a long description that should wrap inside the repository details pane.",
			"Last pushed: Thu, 30 Apr 2026 13:00:00 UTC",
			"Topics: compiler, go, golang",
		} {
			if !strings.Contains(text, want) {
				t.Errorf("width %d: details missing %q; got %q", width, want, text)
			}
		}
	}
}

func TestDetailsBorderSurvivesTabContainer(t *testing.T) {
	for _, width := range []int{60, 79, 80, 120} {
		t.Run(fmt.Sprintf("width %d", width), func(t *testing.T) {
			model := NewModel()
			updated, _ := model.Update(reposMsg{repos: []repo{{name: "card", description: "Repository details"}}})
			manager, err := tabs.NewTabsManager([]tabs.Tab{NewModel(), updated.(Model)})
			if err != nil {
				t.Fatal(err)
			}
			updated, _ = manager.Update(tea.WindowSizeMsg{Width: width, Height: 24})
			updated, _ = updated.Update(tea.KeyPressMsg{Code: tea.KeyTab})
			view := ansi.Strip(updated.View().Content)
			// The outer tab container consumes two border columns and leaves
			// the last terminal column unused. Both detail borders must remain
			// complete on a single line after the outer container renders them.
			contentWidth := width - 3
			paneWidth := contentWidth - contentWidth/2
			border := "+" + strings.Repeat("-", paneWidth-2) + "+"
			if got := strings.Count(view, border); got != 2 {
				t.Errorf("got %d complete detail borders, want 2:\n%s", got, view)
			}
			for i, line := range strings.Split(view, "\n") {
				if got := lipgloss.Width(line); got >= width {
					t.Errorf("line %d uses %d columns, terminal width %d", i, got, width)
				}
			}
		})
	}
}
