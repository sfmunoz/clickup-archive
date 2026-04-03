package tui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ItemDelegate interface
// https://github.com/charmbracelet/bubbles/tree/main/list
// https://pkg.go.dev/github.com/charmbracelet/bubbles/list#ItemDelegate

type itemDelegateMsg struct {
	index int
}

type itemDelegate struct {
	totStyles    int
	itemStyle    [6]lipgloss.Style
	selItemStyle [6]lipgloss.Style
}

func newItemDelegate() *itemDelegate {
	itemStyle := [6]lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}
	selItemStyle := [6]lipgloss.Style{}
	for i, v := range itemStyle {
		selItemStyle[i] = v.Bold(true)
		itemStyle[i] = v.PaddingLeft(2)
	}
	return &itemDelegate{
		totStyles:    len(itemStyle),
		itemStyle:    itemStyle,
		selItemStyle: selItemStyle,
	}
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(_ tea.Msg, l *list.Model) tea.Cmd {
	return func() tea.Msg {
		return itemDelegateMsg{index: l.Index()}
	}
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*item)
	if !ok {
		return
	}
	n := i.GetLevel() % d.totStyles
	pref := ""
	s := d.itemStyle[n]
	if index == m.Index() {
		pref = "> "
		s = d.selItemStyle[n]
	}
	fmt.Fprint(w, s.Render(fmt.Sprintf("%s %2d. %s (%d)", pref, i.GetPos()+1, i.Title(), index+1)))
}
