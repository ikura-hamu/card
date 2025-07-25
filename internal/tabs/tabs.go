package tabs

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
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
	view += fmt.Sprintf("current tab: %s\n", tm.tabNames[tm.activeTab])
	view += tm.tabs[tm.activeTab].View() + "\n"

	return view
}
