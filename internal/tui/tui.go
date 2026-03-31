package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33"))

	levelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

var levelNames = []string{"Workspace", "Space", "Folder", "List", "Task", "Comment"}

type entry struct {
	id   string
	name string
	dir  string
}

type levelState struct {
	dir     string
	entries []entry
	cursor  int
}

type Tui struct {
	clickupDir string
	table      table.Model
	stack      []levelState
	current    []entry
	currentDir string
	level      int
	breadcrumb []string
	err        error
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

func buildTable(entries []entry) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		{Title: "Name", Width: 40},
	}
	rows := make([]table.Row, len(entries))
	for i, e := range entries {
		rows[i] = table.Row{e.id, e.name}
	}
	height := min(len(rows)+1, 20)
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(64),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return t
}

func NewTui(clickupDir string) (*Tui, error) {
	entries, err := loadEntries(clickupDir)
	if err != nil {
		return nil, fmt.Errorf("load workspaces: %w", err)
	}
	return &Tui{
		clickupDir: clickupDir,
		current:    entries,
		currentDir: clickupDir,
		level:      0,
		table:      buildTable(entries),
	}, nil
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) navigateIn() (tea.Model, tea.Cmd) {
	if len(t.current) == 0 || t.level >= len(levelNames)-1 {
		return t, nil
	}
	sel := t.table.Cursor()
	if sel < 0 || sel >= len(t.current) {
		return t, nil
	}
	selectedEntry := t.current[sel]
	nextDir := selectedEntry.dir
	if t.level == 4 { // task → comments
		nextDir = filepath.Join(selectedEntry.dir, "comments")
	}
	entries, err := loadEntries(nextDir)
	if err != nil {
		t.err = fmt.Errorf("cannot open %s: %w", nextDir, err)
		return t, nil
	}
	t.stack = append(t.stack, levelState{
		dir:     t.currentDir,
		entries: t.current,
		cursor:  sel,
	})
	t.breadcrumb = append(t.breadcrumb, selectedEntry.name)
	t.current = entries
	t.currentDir = nextDir
	t.level++
	t.table = buildTable(entries)
	t.err = nil
	return t, nil
}

func (t *Tui) navigateOut() (tea.Model, tea.Cmd) {
	if len(t.stack) == 0 {
		return t, nil
	}
	prev := t.stack[len(t.stack)-1]
	t.stack = t.stack[:len(t.stack)-1]
	t.breadcrumb = t.breadcrumb[:len(t.breadcrumb)-1]
	t.current = prev.entries
	t.currentDir = prev.dir
	t.level--
	t.table = buildTable(prev.entries)
	t.table.SetCursor(prev.cursor)
	t.err = nil
	return t, nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "right", "enter":
			return t.navigateIn()
		case "left", "esc":
			return t.navigateOut()
		}
	}
	var cmd tea.Cmd
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t *Tui) View() tea.View {
	var b strings.Builder
	if len(t.breadcrumb) > 0 {
		b.WriteString(breadcrumbStyle.Render("  "+strings.Join(t.breadcrumb, " › ")) + "\n")
	}
	b.WriteString(levelStyle.Render("  "+levelNames[t.level]) + "\n")
	b.WriteString(baseStyle.Render(t.table.View()) + "\n")
	if t.err != nil {
		b.WriteString(errorStyle.Render("  "+t.err.Error()) + "\n")
	}
	b.WriteString("  " + t.table.HelpView() + "   ←/esc back   →/enter in\n")
	return tea.NewView(b.String())
}

func (t *Tui) Run() error {
	p := tea.NewProgram(t)
	_, err := p.Run()
	return err
}
