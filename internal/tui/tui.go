package tui

import (
	"github.com/sfmunoz/clickup-archive/internal/archive"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

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

type Tui struct {
	archive       *archive.Archive
	items         *Items
	stats         *Stats
	width, height int
}

func NewTui(a *archive.Archive) (*Tui, error) {
	items, err := NewItems(a)
	if err != nil {
		return nil, err
	}
	stats, err := NewStats(a)
	if err != nil {
		return nil, err
	}
	return &Tui{
		archive: a,
		items:   items,
		stats:   stats,
	}, nil
}

func (t *Tui) Init() tea.Cmd {
	return nil
}

func (t *Tui) updateWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	t.width = msg.Width
	t.height = msg.Height
	topbarH := 1
	statusbarH := 1
	sidebarW := 40
	bodyH := t.height - topbarH - statusbarH
	t.items.SetSize(sidebarW, bodyH)
	var cmd tea.Cmd
	t.stats, cmd = t.stats.Update(msg)
	return t, cmd
}

func (t *Tui) updateKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg.String() {
	case "ctrl+c", "q":
		return t, tea.Quit
	case "s":
		var cmd tea.Cmd
		t.stats, cmd = t.stats.Update(StatsVisibleToggleMsg{})
		cmds = append(cmds, cmd)
	default:
		var cmd tea.Cmd
		t.items, cmd = t.items.Update(msg)
		cmds = append(cmds, cmd)
	}
	return t, tea.Batch(cmds...)
}

func (t *Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case itemDelegateMsg:
		return t, nil
	case tea.WindowSizeMsg:
		return t.updateWindowSize(msg)
	case tea.KeyPressMsg:
		return t.updateKeyPress(msg)
	}
	return t, nil
}

func (t *Tui) View() tea.View {
	if t.width == 0 {
		return tea.View{}
	}
	topbarH := 1
	statusbarH := 1
	sidebarW := 40
	contentW := t.width - sidebarW
	bodyH := t.height - topbarH - statusbarH
	topbar := topbarStyle.
		Width(t.width).
		Height(topbarH).
		Render("clickup-archive")
	sidebar := sidebarStyle.
		Width(sidebarW).
		Height(bodyH).
		Render(t.items.View())
	contentInnerW := contentW - sidebarStyle.GetHorizontalFrameSize()
	node := t.items.SelectedNode()
	content := contentStyle.
		Width(contentInnerW).
		Height(bodyH).
		Render(renderContent(node, contentInnerW, bodyH))
	statusbar := statusbarStyle.
		Width(t.width).
		Height(statusbarH).
		Render("q/ctrl-c: quit ; s: show/hide stats")
	screen := lipgloss.JoinVertical(
		lipgloss.Top,
		topbar,
		lipgloss.JoinHorizontal(lipgloss.Top, sidebar, content),
		statusbar,
	)
	statsView := t.stats.View()
	statsW, statsH := lipgloss.Size(statsView)
	comp := lipgloss.NewCompositor(
		lipgloss.NewLayer(screen).X(0).Y(0).Z(0),
		lipgloss.NewLayer(statsView).X((t.width-statsW)/2).Y((t.height-statsH)/2).Z(10),
	)
	var v tea.View
	v.AltScreen = true
	v.SetContent(comp.Render())
	return v
}

func (t *Tui) Run() error {
	p := tea.NewProgram(t)
	_, err := p.Run()
	return err
}
