package tabs

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type CustomHelpTab interface {
	Tab
	Help() help.Model
}

var _ help.KeyMap = (*tabsKeyMap)(nil)

type tabsKeyMap struct {
	left          key.Binding
	right         key.Binding
	quit          key.Binding
	contentKeyMap help.KeyMap
}

func (k tabsKeyMap) ShortHelp() []key.Binding {
	b := []key.Binding{
		k.left, k.right, k.quit,
	}

	if k.contentKeyMap == nil {
		return b
	}

	return append(b, k.contentKeyMap.ShortHelp()...)
}

func (k tabsKeyMap) FullHelp() [][]key.Binding {
	b := [][]key.Binding{
		{k.left, k.right, k.quit},
	}

	if k.contentKeyMap == nil {
		return b
	}

	return append(b, k.contentKeyMap.FullHelp()...)
}

func (tm TabsManager) renderHelp() string {
	if t, ok := tm.tabs[tm.activeTab].(CustomHelpTab); ok {
		return t.Help().View(tm.keyMap)
	}
	return tm.help.View(tm.keyMap)
}
