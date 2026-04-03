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

type itemDelegate struct {
	itemStyle    lipgloss.Style
	selItemStyle lipgloss.Style
}

func newItemDelegate() *itemDelegate {
	return &itemDelegate{
		itemStyle:    lipgloss.NewStyle().PaddingLeft(2),
		selItemStyle: lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("6")),
	}
}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	pref := ""
	fn := d.itemStyle.Render
	if index == m.Index() {
		pref = "> "
		fn = d.selItemStyle.Render
	}
	fmt.Fprint(w, fn(fmt.Sprintf("%s%3d. %s", pref, index+1, i.title)))
}
