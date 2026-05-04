package tui

// https://github.com/charmbracelet/bubbles/tree/main/list
// type Item interface {
//   FilterValue() string
// }
// type DefaultItem interface {
//   Item
//   Title() string
//   Description() string
// }

type item struct {
	node                 any
	title, desc          string
	pos, level           int
	expandable, expanded bool
}

func newItem(node any, title, desc string, pos, level int, expandable, expanded bool) *item {
	return &item{
		node:       node,
		title:      title,
		desc:       desc,
		pos:        pos,
		level:      level,
		expandable: expandable,
		expanded:   expanded,
	}
}

func (i *item) Title() string {
	return i.title
}

func (i *item) Description() string {
	return i.desc
}

func (i *item) FilterValue() string {
	return i.title + " " + i.desc
}

func (i *item) GetPos() int {
	return i.pos
}

func (i *item) GetLevel() int {
	return i.level
}

func (i *item) IsExpandable() bool {
	return i.expandable
}

func (i *item) IsExpanded() bool {
	return i.expanded
}
