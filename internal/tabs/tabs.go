package tabs

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.ikura-hamu.work/card/internal/common/merrors"
	"go.ikura-hamu.work/card/internal/common/size"
)

type Tab interface {
	tea.Model
	Name() string
	KeyMap() help.KeyMap
}

type TabsManager struct {
	tabNames  []string
	activeTab int
	tabs      []Tab
	help      help.Model
	keyMap    tabsKeyMap

	size size.Size
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
		help:      help.New(),
		keyMap: tabsKeyMap{
			left:  key.NewBinding(key.WithKeys("tab", "right"), key.WithHelp("tab", "move left")),
			right: key.NewBinding(key.WithKeys("shift+tab", "left"), key.WithHelp("shift+tab", "move right")),
			quit:  key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
		},
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

func (tm *TabsManager) changeActiveTab(idx int) {
	tm.activeTab = idx
	tm.keyMap.contentKeyMap = tm.tabs[tm.activeTab].KeyMap()
}

func (tm TabsManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			return tm, tea.Quit
		case "tab":
			tm.changeActiveTab((tm.activeTab + 1) % len(tm.tabs))
		case "shift+tab":
			tm.changeActiveTab((tm.activeTab - 1 + len(tm.tabs)) % len(tm.tabs))
		}
	}

	if msg, ok := msg.(error); ok {
		tea.Printf("error: %v", msg)
		return tm, tea.Quit
	}

	// WindowSizeMsg を受け取った場合、コンテンツエリアのサイズを計算して渡す
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		contentWidth, contentHeight := tm.calculateContentSize(msg.Width, msg.Height)
		contentSizeMsg := tea.WindowSizeMsg{
			Width:  contentWidth,
			Height: contentHeight,
		}

		tm.help.Width = contentWidth

		tm.size = size.Size{Width: contentWidth, Height: contentHeight}

		cmds := make([]tea.Cmd, 0, len(tm.tabs))
		for i, tab := range tm.tabs {
			tab, cmd := tab.Update(contentSizeMsg)
			cmds = append(cmds, cmd)
			if t, ok := tab.(Tab); ok {
				tm.tabs[i] = t
			} else {
				return tm, merrors.NewCmd(fmt.Errorf("tab is not a valid Tab type: %T", tab))
			}
		}

		tm.changeActiveTab(tm.activeTab) // 起動時用に呼び出す

		return tm, tea.Batch(cmds...)
	}

	// キーボードのイベントはactive tabにのみ送信
	if msg, ok := msg.(tea.KeyMsg); ok {
		tab, cmd := tm.tabs[tm.activeTab].Update(msg)
		if t, ok := tab.(Tab); ok {
			tm.tabs[tm.activeTab] = t
		} else {
			return tm, merrors.NewCmd(fmt.Errorf("tab is not a valid Tab type: %T", tab))
		}
		return tm, cmd
	}

	cmds := make([]tea.Cmd, 0, len(tm.tabs))
	for i, tab := range tm.tabs {
		tab, cmd := tab.Update(msg)
		cmds = append(cmds, cmd)
		if t, ok := tab.(Tab); ok {
			tm.tabs[i] = t
		} else {
			return tm, merrors.NewCmd(fmt.Errorf("tab is not a valid Tab type: %T", tab))
		}
	}

	return tm, tea.Batch(cmds...)
}

func (tm TabsManager) renderTabHeaders() string {
	tabHeaders := make([]string, 0, len(tm.tabs))
	for i, name := range tm.tabNames {
		if i == tm.activeTab {
			tabHeaders = append(tabHeaders, activeTabStyle.Render(fmt.Sprintf(" %s ", name)))
		} else {
			tabHeaders = append(tabHeaders, inactiveTabStyle.Render(fmt.Sprintf(" %s ", name)))
		}
	}

	tabHeader := tabHeaderStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, tabHeaders...)) + "\n"

	return tabHeader
}

// calculateContentSize calculates the available content area size
func (tm TabsManager) calculateContentSize(terminalWidth, terminalHeight int) (width int, height int) {
	// タブヘッダーのサンプルを作成してサイズを測定
	sampleTabHeader := tm.renderTabHeaders()
	tabHeaderHeight := lipgloss.Height(sampleTabHeader)

	tabHelp := tm.renderHelp()
	tabHelpHeight := lipgloss.Height(tabHelp)

	// コンテンツスタイルのボーダーとパディングを考慮
	contentSample := contentStyle.Render("sample")
	contentBorderWidth := lipgloss.Width(contentSample) - lipgloss.Width("sample")
	contentBorderHeight := lipgloss.Height(contentSample) - lipgloss.Height("sample")

	// 利用可能なコンテンツエリアのサイズを計算
	contentWidth := terminalWidth - contentBorderWidth
	contentHeight := terminalHeight - tabHeaderHeight - tabHelpHeight - contentBorderHeight

	// 最小サイズを確保
	if contentWidth < 1 {
		contentWidth = 1
	}
	if contentHeight < 1 {
		contentHeight = 1
	}

	return contentWidth, contentHeight
}

func (tm TabsManager) View() string {
	var view string

	view += tm.renderTabHeaders()
	view += contentStyle.Width(tm.size.Width).Height(tm.size.Height).
		Render(tm.tabs[tm.activeTab].View())
	view += "\n" + tm.renderHelp()

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
