package tui

import (
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2/tree"
)

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
			return "[" + n.Name + "]"
		}
		return n.Name
	}
	if len(n.Children) < 1 {
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

func (n *Node) Init() tea.Cmd {
	return nil
}

func (n *Node) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "s":
			n.SetCursor(!n.Cursor)
			return n, nil
		}
	}
	return n, nil
}

func (n *Node) View() tea.View {
	var b strings.Builder
	b.WriteString(n.BuildTree().String() + "\n")
	return tea.NewView(b.String())
}

func TreeDemo() *Node {
	return NewNode("ClickUp Archive").SetOpen(true).SetCursor(true).SetChildren(
		NewNode("workspace-1").SetOpen(true).SetChildren(
			NewNode("space-1"),
			NewNode("space-2").SetOpen(true).SetChildren(
				NewNode("folder-21").SetOpen(true).SetChildren(
					NewNode("list-211"),
					NewNode("list-212"),
				),
				NewNode("folder-22").SetOpen(false).SetChildren(
					NewNode("list-221"),
					NewNode("list-222"),
				),
			),
			NewNode("space-3").SetOpen(true).SetChildren(
				NewNode("folder-31").SetOpen(true).SetChildren(
					NewNode("list-311"),
					NewNode("list-312"),
				),
			),
		),
	)
}
