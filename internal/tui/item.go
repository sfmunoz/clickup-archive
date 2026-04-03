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
	title, desc string
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.desc
}

func (i item) FilterValue() string {
	return i.title
}
