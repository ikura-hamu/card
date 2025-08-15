package merrors

import (
	tea "github.com/charmbracelet/bubbletea"
)

func NewCmd(err error) func() tea.Msg {
	return func() tea.Msg {
		return err
	}
}
