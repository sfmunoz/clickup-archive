package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

type Stats struct {
	archive *archive.Archive
	visible bool
}

type StatsVisibleToggleMsg struct{}

func NewStats(a *archive.Archive) (*Stats, error) {
	return &Stats{
		archive: a,
		visible: false,
	}, nil
}

func (s *Stats) Update(msg tea.Msg) (*Stats, tea.Cmd) {
	switch msg.(type) {
	case StatsVisibleToggleMsg:
		s.visible = !s.visible
	}
	return s, nil
}

func (s *Stats) View() string {
	if !s.visible {
		return ""
	}
	var wTot, sTot, fTot, lTot, tTot, cTot = 0, 0, 0, 0, 0, 0
	for _, v1 := range s.archive.Children {
		wTot++
		for _, v2 := range v1.Children {
			sTot++
			for _, v3 := range v2.Children {
				fTot++
				for _, v4 := range v3.Children {
					lTot++
					for _, v5 := range v4.Children {
						tTot++
						for range v5.Children {
							cTot++
						}
					}
				}
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "workspaces ... %4d\n", wTot)
	fmt.Fprintf(&sb, "spaces ....... %4d\n", sTot)
	fmt.Fprintf(&sb, "folders ...... %4d\n", fTot)
	fmt.Fprintf(&sb, "lists ........ %4d\n", lTot)
	fmt.Fprintf(&sb, "tasks ........ %4d\n", tTot)
	fmt.Fprintf(&sb, "comments ..... %4d\n", cTot)
	return sb.String()
}
