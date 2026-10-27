package repo

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"go.ikura-hamu.work/card/internal/common/theme"
)

func themedDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(lipgloss.Color(theme.Text))
	d.Styles.NormalDesc = d.Styles.NormalDesc.Foreground(lipgloss.Color(theme.Secondary))
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(lipgloss.Color(theme.Primary)).
		BorderForeground(lipgloss.Color(theme.Primary)).Bold(true)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(lipgloss.Color(theme.Secondary)).
		BorderForeground(lipgloss.Color(theme.Primary))
	d.Styles.DimmedTitle = d.Styles.DimmedTitle.Foreground(lipgloss.Color(theme.Secondary))
	d.Styles.DimmedDesc = d.Styles.DimmedDesc.Foreground(lipgloss.Color(theme.Secondary))
	d.Styles.FilterMatch = d.Styles.FilterMatch.Foreground(lipgloss.Color(theme.Accent))
	return d
}

func styleRepositoryList(m *list.Model) {
	primary := lipgloss.Color(theme.Primary)
	secondary := lipgloss.Color(theme.Secondary)
	accent := lipgloss.Color(theme.Accent)
	m.Styles.Title = m.Styles.Title.Background(primary).Foreground(lipgloss.Color(theme.Inverse))
	m.Styles.StatusBar = m.Styles.StatusBar.Foreground(secondary)
	m.Styles.StatusEmpty = m.Styles.StatusEmpty.Foreground(secondary)
	m.Styles.StatusBarActiveFilter = m.Styles.StatusBarActiveFilter.Foreground(accent)
	m.Styles.StatusBarFilterCount = m.Styles.StatusBarFilterCount.Foreground(secondary)
	m.Styles.NoItems = m.Styles.NoItems.Foreground(secondary)
	m.Styles.DividerDot = m.Styles.DividerDot.Foreground(primary)
	m.Styles.ArabicPagination = m.Styles.ArabicPagination.Foreground(secondary)
	m.Styles.ActivePaginationDot = m.Styles.ActivePaginationDot.Foreground(accent)
	m.Styles.InactivePaginationDot = m.Styles.InactivePaginationDot.Foreground(secondary)
	m.Paginator.ActiveDot = m.Styles.ActivePaginationDot.String()
	m.Paginator.InactiveDot = m.Styles.InactivePaginationDot.String()
	m.Styles.Filter.Focused.Prompt = m.Styles.Filter.Focused.Prompt.Foreground(primary)
	m.Styles.Filter.Focused.Text = m.Styles.Filter.Focused.Text.Foreground(lipgloss.Color(theme.Text))
	m.Styles.Filter.Focused.Placeholder = m.Styles.Filter.Focused.Placeholder.Foreground(secondary)
	m.Styles.Filter.Blurred = m.Styles.Filter.Focused
	m.Styles.Filter.Cursor.Color = accent
	m.FilterInput.SetStyles(m.Styles.Filter)
	m.Help.Styles = theme.HelpStyles()
}
