package tui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Items struct {
	archive *archive.Archive
	list    list.Model
	width   int
	height  int
}

func NewItems(a *archive.Archive) (*Items, error) {
	items := make([]list.Item, 0)
	for _, v1 := range a.Children {
		items = append(items, item{title: ".." + v1.Data.ID, desc: v1.Data.Name})
		for _, v2 := range v1.Children {
			items = append(items, item{title: "...." + v2.Data.ID, desc: v2.Data.Name})
			for _, v3 := range v2.Children {
				items = append(items, item{title: "......" + v3.Data.ID, desc: v3.Data.Name})
				for _, v4 := range v3.Children {
					items = append(items, item{title: "........" + v4.Data.ID, desc: v4.Data.Name})
					for _, v5 := range v4.Children {
						items = append(items, item{title: ".........." + v5.Data.ID, desc: v5.Data.Name})
						for _, v6 := range v5.Children {
							items = append(items, item{title: "............" + v6.Data.ID, desc: v6.Data.Text})
						}
					}
				}
			}
		}
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Item list"
	return &Items{
		archive: a,
		list:    list,
		width:   0,
		height:  0,
	}, nil
}

func (i *Items) Update(msg tea.Msg) (*Items, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.width = msg.Width
		i.height = msg.Height
		i.list.SetWidth(i.width - 20)   // FIXME
		i.list.SetHeight(i.height - 10) // FIXME
		var cmd tea.Cmd
		i.list, cmd = i.list.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyPressMsg:
		var cmd tea.Cmd
		i.list, cmd = i.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	return i, tea.Batch(cmds...)
}

func (i *Items) View() string {
	if i.width == 0 {
		return ""
	}
	return docStyle.Render(i.list.View())
}
