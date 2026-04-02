package tui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/archive"

	tea "charm.land/bubbletea/v2"
)

type entry struct {
	id   string
	name string
	dir  string
}

type Tui struct {
	archive *archive.Archive
	tree    *Node
}

func loadEntries(dir string) ([]entry, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var entries []entry
	for _, de := range dirEntries {
		if !de.IsDir() {
			continue
		}
		indexPath := filepath.Join(dir, de.Name(), "index.json")
		data, err := os.ReadFile(indexPath)
		if err != nil {
			continue
		}
		var v struct {
			Name string `json:"name"`
			Text string `json:"comment_text"`
		}
		if err := json.Unmarshal(data, &v); err != nil {
			continue
		}
		name := v.Name
		if name == "" {
			name = v.Text
		}
		entries = append(entries, entry{
			id:   de.Name(),
			name: name,
			dir:  filepath.Join(dir, de.Name()),
		})
	}
	return entries, nil
}

func (t *Tui) loadNodeChildren(n *Node) {
	if n.childrenLoaded || n.childrenDir == "" {
		return
	}
	n.childrenLoaded = true
	entries, err := loadEntries(n.childrenDir)
	if err != nil {
		return
	}
	for _, e := range entries {
		child := &Node{
			Name:  e.name,
			dir:   e.dir,
			level: n.level + 1,
		}
		if child.level == 4 { // task → children live in comments/
			child.childrenDir = filepath.Join(e.dir, "comments")
		} else if child.level < 5 { // comment nodes have no children
			child.childrenDir = e.dir
		}
		n.AppendChild(child)
	}
}

func NewTui(a *archive.Archive) (*Tui, error) {
	root := &Node{
		Name:        "ClickUp Archive",
		dir:         a.GetDir(),
		childrenDir: a.GetDir(),
		level:       -1,
		Open:        true,
		Cursor:      true,
	}
	t := &Tui{archive: a, tree: root}
	t.loadNodeChildren(root)

	return t, nil
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if kp, ok := msg.(tea.KeyPressMsg); ok {
		switch kp.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "right", "enter", " ":
			cursor := t.tree.root().findCursor()
			if cursor != nil {
				t.loadNodeChildren(cursor)
			}
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
