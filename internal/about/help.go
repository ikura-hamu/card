package about

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

type KeyMap viewport.KeyMap

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(
			key.WithKeys(k.Up.Keys()...),
			key.WithHelp(k.Up.Help().Key, k.Up.Help().Desc),
		),
		key.NewBinding(
			key.WithKeys(k.Down.Keys()...),
			key.WithHelp(k.Down.Help().Key, k.Down.Help().Desc),
		),
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.NewBinding(
				key.WithKeys(k.Up.Keys()...),
				key.WithHelp(k.Up.Help().Key, k.Up.Help().Desc),
			),
		},
		{
			key.NewBinding(
				key.WithKeys(k.Down.Keys()...),
				key.WithHelp(k.Down.Help().Key, k.Down.Help().Desc),
			),
		},
		{
			key.NewBinding(
				key.WithKeys(k.PageUp.Keys()...),
				key.WithHelp(k.PageUp.Help().Key, k.PageUp.Help().Desc),
			),
		},
		{
			key.NewBinding(
				key.WithKeys(k.PageDown.Keys()...),
				key.WithHelp(k.PageDown.Help().Key, k.PageDown.Help().Desc),
			),
		},
	}
}
