package tabs

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.ikura-hamu.work/card/internal/common/merrors"
)

type Tab interface {
	tea.Model
	Name() string
}

type TabsManager struct {
	tabNames  []string
	activeTab int
	tabs      []Tab
}

func NewTabsManager(tabs []Tab) (TabsManager, error) {
	if len(tabs) == 0 {
		return TabsManager{}, fmt.Errorf("at least one tab is required")
	}

	tabNames := make([]string, len(tabs))
	for i, tab := range tabs {
		tabNames[i] = tab.Name()
	}
	return TabsManager{
		tabNames:  tabNames,
		activeTab: 0,
		tabs:      tabs,
	}, nil
}

func (tm TabsManager) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(tm.tabs))
	for _, tab := range tm.tabs {
		cmd := tab.Init()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (tm TabsManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "q":
			return tm, tea.Quit
		case "tab":
			tm.activeTab++
			tm.activeTab %= len(tm.tabs)
		case "shift+tab":
			tm.activeTab--
			if tm.activeTab < 0 {
				tm.activeTab += len(tm.tabs)
			}
			tm.activeTab %= len(tm.tabs)
		}
	}

	if msg, ok := msg.(merrors.Msg); ok {
		log.Printf("error: %v", msg.Error())
		return tm, tea.Quit
	}

	// キーボードのイベントはactive tabにのみ送信
	if msg, ok := msg.(tea.KeyMsg); ok {
		tab, cmd := tm.tabs[tm.activeTab].Update(msg)
		if tab, ok := tab.(Tab); ok {
			tm.tabs[tm.activeTab] = tab
		} else {
			return tm, tea.Quit // If a tab is not a valid Tab type, exit
		}
		return tm, cmd
	}

	cmds := make([]tea.Cmd, 0, len(tm.tabs))
	for i, tab := range tm.tabs {
		tab, cmd := tab.Update(msg)
		cmds = append(cmds, cmd)
		if tab, ok := tab.(Tab); ok {
			tm.tabs[i] = tab
		} else {
			return tm, tea.Quit // If a tab is not a valid Tab type, exit
		}
	}

	return tm, tea.Batch(cmds...)
}

func (tm TabsManager) View() string {
	var view string

	tabHeaders := make([]string, 0, len(tm.tabs))
	for i, name := range tm.tabNames {
		if i == tm.activeTab {
			tabHeaders = append(tabHeaders, activeTabStyle.Render(fmt.Sprintf(" %s ", name)))
		} else {
			tabHeaders = append(tabHeaders, inactiveTabStyle.Render(fmt.Sprintf(" %s ", name)))
		}
	}

	tabHeader := tabHeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, tabHeaders...)) + "\n"

	view += tabHeader
	view += contentStyle.Render(tm.tabs[tm.activeTab].View())

	return view
}

var (
	activeTabStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, true, false, true).
			BorderForeground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("240")).
			Foreground(lipgloss.Color("white")).
			Margin(0, 1).
			Align(lipgloss.Center)
	inactiveTabStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder(), true, true, false, true).
				BorderForeground(lipgloss.Color("240")).
				Background(lipgloss.Color("white")).
				Foreground(lipgloss.Color("black")).
				Margin(0, 1).
				Align(lipgloss.Center)
	tabHeaderStyle = lipgloss.NewStyle().
			Margin(0, 1)
	contentStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Background(lipgloss.Color("white")).
			Foreground(lipgloss.Color("black")).
			Align(lipgloss.Left)
)
