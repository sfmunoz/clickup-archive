package tui

import (
	"github.com/sfmunoz/clickup-archive/internal/archive"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Tui struct {
	archive *archive.Archive
	width   int
	height  int
}

func NewTui(a *archive.Archive) (*Tui, error) {
	return &Tui{archive: a}, nil
}

var (
	topbarStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#5C5C5C")).
			PaddingLeft(1)

	sidebarStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#3C3C3C"))

	contentStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	statusbarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Background(lipgloss.Color("#2A2A2A")).
			PaddingLeft(1)
)

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return t, tea.Quit
		}
	}
	return t, nil
}

func (t *Tui) View() tea.View {
	if t.width == 0 {
		return tea.View{}
	}
	topbarH := 1
	statusbarH := 1
	sidebarW := 25
	contentW := t.width - sidebarW
	bodyH := t.height - topbarH - statusbarH
	topbar := topbarStyle.
		Width(t.width).
		Height(topbarH).
		Render("clickup-archive")
	sidebar := sidebarStyle.
		Width(sidebarW).
		Height(bodyH).
		Render("Workspace\n  Space A\n  Space B\n    Folder 1\n    Folder 2\n    Folder 3")
	content := contentStyle.
		Width(contentW - lipgloss.Width(sidebarStyle.Render(""))). // subtract border
		Height(bodyH).
		Render("Select a list to browse tasks.")
	statusbar := statusbarStyle.
		Width(t.width).
		Height(statusbarH).
		Render("q/ctrl-c: quit")
	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content)
	screen := lipgloss.JoinVertical(lipgloss.Top, topbar, body, statusbar)
	var v tea.View
	v.AltScreen = true
	v.SetContent(screen)
	return v
}

func (t *Tui) Run() error {
	p := tea.NewProgram(t)
	_, err := p.Run()
	return err
}
