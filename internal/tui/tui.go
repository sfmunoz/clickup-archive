package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func tableBuild() table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Name", Width: 30},
	}
	rows := []table.Row{
		{"1", "Workspace"},
		{"2", "Space"},
		{"3", "Folder"},
		{"4", "List"},
		{"5", "Task"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
		table.WithWidth(42),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("240")).
		Bold(false)
	t.SetStyles(s)
	return t
}

type Tui struct {
	clickupDir string
	table      table.Model
}

func NewTui(clickupDir string) (*Tui, error) {
	return &Tui{
		clickupDir: clickupDir,
		table:      tableBuild(),
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
		case "enter", "space":
			return t, tea.Batch(
				tea.Printf("selection: '%s'", t.table.SelectedRow()[1]),
			)
		}
	}
	var cmd tea.Cmd
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

func (t *Tui) View() tea.View {
	return tea.NewView(baseStyle.Render(t.table.View()) + "\n  " + t.table.HelpView() + "\n")
}

func (t *Tui) Run() error {
	p := tea.NewProgram(t)
	_, err := p.Run()
	return err
}
