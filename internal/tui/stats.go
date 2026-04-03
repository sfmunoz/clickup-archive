package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

var statsStyle = lipgloss.NewStyle().
	AlignHorizontal(lipgloss.Center).
	Border(lipgloss.RoundedBorder(), true).
	Padding(1, 5)

var statsTitleStyle = lipgloss.NewStyle().Bold(true).Underline(true)

type Stats struct {
	archive *archive.Archive
}

func NewStats(a *archive.Archive) (*Stats, error) {
	return &Stats{
		archive: a,
	}, nil
}

func (s *Stats) View() string {
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
	fmt.Fprint(&sb, statsTitleStyle.Render("stats")+"\n\n")
	fmt.Fprintf(&sb, "workspaces ... %4d\n", wTot)
	fmt.Fprintf(&sb, "spaces ....... %4d\n", sTot)
	fmt.Fprintf(&sb, "folders ...... %4d\n", fTot)
	fmt.Fprintf(&sb, "lists ........ %4d\n", lTot)
	fmt.Fprintf(&sb, "tasks ........ %4d\n", tTot)
	fmt.Fprintf(&sb, "comments ..... %4d", cTot)
	return statsStyle.Render(sb.String())
}
