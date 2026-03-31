package tui

import (
	"fmt"

	"github.com/sfmunoz/logit"

	tea "charm.land/bubbletea/v2"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Tui struct {
	clickupDir string
	choices    []string
	selected   map[int]struct{}
	cursor     int
}

func NewTui(clickupDir string) (*Tui, error) {
	return &Tui{
		clickupDir: clickupDir,
		choices:    []string{"Workspace", "Space", "Folder", "List", "Task"},
		selected:   make(map[int]struct{}),
		cursor:     0,
	}, nil
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}
		case "down", "j":
			if t.cursor < len(t.choices)-1 {
				t.cursor++
			}
		case "enter", "space":
			_, ok := t.selected[t.cursor]
			if ok {
				delete(t.selected, t.cursor)
			} else {
				t.selected[t.cursor] = struct{}{}
			}
		}
	}
	return t, nil
}

func (t *Tui) View() tea.View {
	s := "\nSelect components:\n\n"
	for i, choice := range t.choices {
		cursor := " "
		if t.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := t.selected[i]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\nPress q to quit.\n"
	v := tea.NewView(s)
	v.WindowTitle = "ClickUp Archive"
	return v
}

func (t *Tui) Run() error {
	p := tea.NewProgram(t)
	_, err := p.Run()
	return err
}
