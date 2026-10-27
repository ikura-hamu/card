package merrors

import (
	tea "charm.land/bubbletea/v2"
)

func NewCmd(err error) func() tea.Msg {
	return func() tea.Msg {
		return err
	}
}
