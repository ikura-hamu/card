package tabs

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
	"go.ikura-hamu.work/card/internal/common/theme"
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
	helpModel := tm.help
	if t, ok := tm.tabs[tm.activeTab].(CustomHelpTab); ok {
		helpModel = t.Help()
	}
	helpModel.Styles = theme.HelpStyles()
	helpModel.SetWidth(tm.size.Width)
	return lipgloss.NewStyle().MaxWidth(tm.size.Width).Render(helpModel.View(tm.keyMap))
}
