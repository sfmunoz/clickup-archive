package tui

import (
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/archive"

	tea "charm.land/bubbletea/v2"
)

type Tui struct {
	archive *archive.Archive
	tree    *Node
}

func buildTree(a *archive.Archive) *Node {
	root := &Node{Name: "ClickUp Archive", Open: true, Cursor: true}
	for _, w := range a.Children {
		wn := &Node{Name: w.Data.Name}
		for _, s := range w.Children {
			sn := &Node{Name: s.Data.Name}
			for _, f := range s.Children {
				fn := &Node{Name: f.Data.Name}
				for _, l := range f.Children {
					ln := &Node{Name: l.Data.Name}
					for _, t := range l.Children {
						tn := &Node{Name: t.Data.Name}
						for _, c := range t.Children {
							tn.AppendChild(&Node{Name: c.Data.Text})
						}
						ln.AppendChild(tn)
					}
					fn.AppendChild(ln)
				}
				sn.AppendChild(fn)
			}
			wn.AppendChild(sn)
		}
		root.AppendChild(wn)
	}
	return root
}

func NewTui(a *archive.Archive) (*Tui, error) {
	return &Tui{archive: a, tree: buildTree(a)}, nil
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if kp, ok := msg.(tea.KeyPressMsg); ok {
		switch kp.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "s":
			return screen{}, tea.RequestWindowSize
		}
	}
	model, cmd := t.tree.Update(msg)
	t.tree = model.(*Node)
	return t, cmd
}

func (t *Tui) View() tea.View {
	var b strings.Builder
	b.WriteString(t.tree.View().Content)
	b.WriteString("  ↑↓ move   ←/→ collapse/expand   q quit\n")
	return tea.NewView(b.String())
}

func (t *Tui) Run() error {
	p := tea.NewProgram(t)
	_, err := p.Run()
	return err
}
