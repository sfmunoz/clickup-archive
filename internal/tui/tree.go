package tui

import (
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
)

var selStyle = lipgloss.NewStyle().
	Background(lipgloss.BrightBlack).
	Padding(0, 1, 0)

type Node struct {
	Name     string
	Parent   *Node
	Children []*Node
	Open     bool
	Cursor   bool
}

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Parent:   nil,
		Children: make([]*Node, 0),
		Open:     false,
		Cursor:   false,
	}
}

func (n *Node) SetName(name string) *Node {
	n.Name = name
	return n
}

func (n *Node) SetParent(parent *Node) *Node {
	n.Parent = parent
	return n
}

func (n *Node) AppendChild(c *Node) *Node {
	if c == nil {
		return n
	}
	if slices.Contains(n.Children, c) {
		return n
	}
	n.Children = append(n.Children, c)
	c.Parent = n
	return n
}

func (n *Node) SetChildren(children ...*Node) *Node {
	n.Children = children
	for _, c := range n.Children {
		c.Parent = n
	}
	return n
}

func (n *Node) SetOpen(open bool) *Node {
	n.Open = open
	return n
}

func (n *Node) SetCursor(cursor bool) *Node {
	n.Cursor = cursor
	return n
}

func (n *Node) String() string {
	name := func() string {
		if n.Cursor {
			return selStyle.Render(n.Name)
		}
		return n.Name
	}
	hasOrMayHaveChildren := len(n.Children) > 0
	if !hasOrMayHaveChildren {
		return name()
	}
	if n.Open {
		return "▼ " + name()
	}
	return "▶ " + name()
}

func (n *Node) BuildTree() *tree.Tree {
	t := tree.Root(n)
	if n.Parent == nil {
		t = t.Enumerator(tree.RoundedEnumerator)
	}
	if n.Open {
		for _, c := range n.Children {
			t.Child(c.BuildTree())
		}
	}
	return t
}

func (n *Node) root() *Node {
	cur := n
	for cur.Parent != nil {
		cur = cur.Parent
	}
	return cur
}

func (n *Node) visibleNodes() []*Node {
	result := []*Node{n}
	if n.Open {
		for _, c := range n.Children {
			result = append(result, c.visibleNodes()...)
		}
	}
	return result
}

func (n *Node) findCursor() *Node {
	if n.Cursor {
		return n
	}
	for _, c := range n.Children {
		if found := c.findCursor(); found != nil {
			return found
		}
	}
	return nil
}

func (n *Node) moveCursorUp() {
	visible := n.root().visibleNodes()
	for i, node := range visible {
		if node.Cursor && i > 0 {
			node.Cursor = false
			visible[i-1].Cursor = true
			return
		}
	}
}

func (n *Node) moveCursorDown() {
	visible := n.root().visibleNodes()
	for i, node := range visible {
		if node.Cursor && i < len(visible)-1 {
			node.Cursor = false
			visible[i+1].Cursor = true
			return
		}
	}
}

func (n *Node) moveLeft() {
	cursor := n.root().findCursor()
	if cursor == nil {
		return
	}
	if cursor.Open {
		cursor.Open = false
	} else if cursor.Parent != nil {
		cursor.Cursor = false
		cursor.Parent.Cursor = true
	}
}

func (n *Node) moveRight() {
	cursor := n.root().findCursor()
	if cursor == nil || len(cursor.Children) == 0 {
		return
	}
	if !cursor.Open {
		cursor.Open = true
	} else {
		cursor.Cursor = false
		cursor.Children[0].Cursor = true
	}
}

func (n *Node) toggleOpen() {
	cursor := n.root().findCursor()
	if cursor == nil || len(cursor.Children) == 0 {
		return
	}
	cursor.Open = !cursor.Open
}

func (n *Node) Init() tea.Cmd {
	return nil
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			n.moveCursorUp()
		case "down":
			n.moveCursorDown()
		case "left":
			n.moveLeft()
		case "right":
			n.moveRight()
		case "enter", " ":
			n.toggleOpen()
		}
	}
	return n, nil
}

func (n *Node) View() tea.View {
	var b strings.Builder
	b.WriteString(n.BuildTree().String() + "\n")
	return tea.NewView(b.String())
}
