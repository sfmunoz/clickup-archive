package tui

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

var itemsStyle = lipgloss.NewStyle().Margin(1, 2)

type Items struct {
	archive  *archive.Archive
	list     list.Model
	expanded map[any]bool
	width    int
	height   int
}

func NewItems(a *archive.Archive) (*Items, error) {
	expanded := make(map[any]bool)
	for _, w := range a.Children {
		expanded[w] = true
	}
	items := buildVisibleItems(a, expanded)
	// delegate := list.NewDefaultDelegate()
	// delegate.SetSpacing(0)
	list := list.New(items, newItemDelegate(), 0, 0)
	list.Title = "Item list"
	return &Items{
		archive:  a,
		list:     list,
		expanded: expanded,
		width:    0,
		height:   0,
	}, nil
}

func (i *Items) SetSize(w, h int) {
	i.width = w
	i.height = h
	i.list.SetWidth(w - itemsStyle.GetHorizontalFrameSize())
	i.list.SetHeight(h - itemsStyle.GetVerticalFrameSize())
}

func (i *Items) Update(msg tea.Msg) (*Items, tea.Cmd) {
	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "+":
			return i, i.expandSelected()
		case "-":
			return i, i.collapseSelected()
		}
	}
	var cmd tea.Cmd
	i.list, cmd = i.list.Update(msg)
	return i, cmd
}

func (i *Items) IsFiltering() bool {
	return i.list.FilterState() == list.Filtering
}

func (i *Items) SelectedNode() any {
	selected := i.list.SelectedItem()
	if it, ok := selected.(*item); ok {
		return it.node
	}
	return nil
}

func (i *Items) expandSelected() tea.Cmd {
	selected := i.SelectedNode()
	if selected == nil || !hasChildren(selected) {
		return nil
	}
	if !i.expanded[selected] {
		i.expanded[selected] = true
		return i.refresh(selected)
	}
	depth := minCollapsedDepth(selected, i.expanded, 0)
	if depth < 0 {
		return nil
	}
	expandCollapsedAtDepth(selected, i.expanded, 0, depth)
	return i.refresh(selected)
}

func (i *Items) collapseSelected() tea.Cmd {
	selected := i.SelectedNode()
	if selected == nil || !hasChildren(selected) {
		return nil
	}
	depth := maxExpandedDescendantDepth(selected, i.expanded, 0)
	if depth < 0 {
		if !i.expanded[selected] {
			return nil
		}
		i.expanded[selected] = false
		return i.refresh(selected)
	}
	collapseExpandedAtDepth(selected, i.expanded, 0, depth)
	return i.refresh(selected)
}

func (i *Items) refresh(selected any) tea.Cmd {
	cmd := i.list.SetItems(buildVisibleItems(i.archive, i.expanded))
	i.selectNode(selected)
	return cmd
}

func (i *Items) selectNode(node any) {
	if node == nil {
		return
	}
	for idx, it := range i.list.Items() {
		item, ok := it.(*item)
		if ok && item.node == node {
			i.list.Select(idx)
			return
		}
	}
}

func (i *Items) View() string {
	if i.width == 0 {
		return ""
	}
	return itemsStyle.Render(i.list.View())
}

type childNode struct {
	node        any
	title, desc string
	pos         int
}

func buildVisibleItems(a *archive.Archive, expanded map[any]bool) []list.Item {
	items := make([]list.Item, 0)
	for _, child := range archiveChildren(a) {
		appendVisibleItem(&items, child, 0, expanded)
	}
	return items
}

func appendVisibleItem(items *[]list.Item, child childNode, level int, expanded map[any]bool) {
	isExpandable := hasChildren(child.node)
	isExpanded := isExpandable && expanded[child.node]
	*items = append(*items, newItem(child.node, child.title, child.desc, child.pos, level, isExpandable, isExpanded))
	if !isExpanded {
		return
	}
	for _, sub := range nodeChildren(child.node) {
		appendVisibleItem(items, sub, level+1, expanded)
	}
}

func archiveChildren(a *archive.Archive) []childNode {
	children := make([]childNode, 0, len(a.Children))
	for pos, w := range a.Children {
		children = append(children, childNode{node: w, title: w.Data.Name, desc: w.Data.ID, pos: pos})
	}
	return children
}

func nodeChildren(node any) []childNode {
	switch n := node.(type) {
	case *archive.Workspace:
		children := make([]childNode, 0, len(n.Children))
		for pos, s := range n.Children {
			children = append(children, childNode{node: s, title: s.Data.Name, desc: s.Data.ID, pos: pos})
		}
		return children
	case *archive.Space:
		children := make([]childNode, 0, len(n.Children))
		for pos, f := range n.Children {
			children = append(children, childNode{node: f, title: f.Data.Name, desc: f.Data.ID, pos: pos})
		}
		return children
	case *archive.Folder:
		children := make([]childNode, 0, len(n.Children))
		for pos, l := range n.Children {
			children = append(children, childNode{node: l, title: l.Data.Name, desc: l.Data.ID, pos: pos})
		}
		return children
	case *archive.List:
		children := make([]childNode, 0, len(n.Children))
		for pos, t := range n.Children {
			children = append(children, childNode{node: t, title: t.Data.ID, desc: t.Data.Name, pos: pos})
		}
		return children
	case *archive.Task:
		children := make([]childNode, 0, len(n.Children))
		for pos, c := range n.Children {
			children = append(children, childNode{node: c, title: c.Data.ID, desc: c.Data.Text, pos: pos})
		}
		return children
	default:
		return nil
	}
}

func hasChildren(node any) bool {
	return len(nodeChildren(node)) > 0
}

func minCollapsedDepth(node any, expanded map[any]bool, depth int) int {
	if hasChildren(node) && !expanded[node] {
		return depth
	}
	minDepth := -1
	for _, child := range nodeChildren(node) {
		childDepth := minCollapsedDepth(child.node, expanded, depth+1)
		if childDepth >= 0 && (minDepth < 0 || childDepth < minDepth) {
			minDepth = childDepth
		}
	}
	return minDepth
}

func expandCollapsedAtDepth(node any, expanded map[any]bool, depth, targetDepth int) {
	if depth == targetDepth {
		if hasChildren(node) {
			expanded[node] = true
		}
		return
	}
	for _, child := range nodeChildren(node) {
		expandCollapsedAtDepth(child.node, expanded, depth+1, targetDepth)
	}
}

func maxExpandedDescendantDepth(node any, expanded map[any]bool, depth int) int {
	maxDepth := -1
	for _, child := range nodeChildren(node) {
		if hasChildren(child.node) && expanded[child.node] && depth+1 > maxDepth {
			maxDepth = depth + 1
		}
		childDepth := maxExpandedDescendantDepth(child.node, expanded, depth+1)
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}
	return maxDepth
}

func collapseExpandedAtDepth(node any, expanded map[any]bool, depth, targetDepth int) {
	if depth == targetDepth {
		if hasChildren(node) {
			expanded[node] = false
		}
		return
	}
	for _, child := range nodeChildren(node) {
		collapseExpandedAtDepth(child.node, expanded, depth+1, targetDepth)
	}
}
