package stats

import (
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/archive"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Stats struct {
	archive *archive.Archive
}

func NewStats(a *archive.Archive) (*Stats, error) {
	return &Stats{archive: a}, nil
}

func (s *Stats) Run() error {
	var (
		wTot int
		sTot int
		fTot int
		lTot int
		tTot int
		cTot int
	)
	for i1, v1 := range s.archive.Children {
		wTot++
		log.Info("..", "i", i1+1, "t", wTot, "id", v1.Data.ID, "n", v1.Data.Name)
		for i2, v2 := range v1.Children {
			sTot++
			log.Info("....", "i", i2+1, "t", sTot, "id", v2.Data.ID, "n", v2.Data.Name)
			for i3, v3 := range v2.Children {
				fTot++
				log.Info("......", "i", i3+1, "t", fTot, "id", v3.Data.ID, "n", v3.Data.Name)
				for i4, v4 := range v3.Children {
					lTot++
					log.Info("........", "i", i4+1, "t", lTot, "id", v4.Data.ID, "n", v4.Data.Name)
					for i5, v5 := range v4.Children {
						tTot++
						log.Info("..........", "i", i5+1, "t", tTot, "id", v5.Data.ID, "n", v5.Data.Name)
						for i6, v6 := range v5.Children {
							cTot++
							log.Info("............", "i", i6+1, "t", cTot, "id", v6.Data.ID, "c", strings.ReplaceAll(
								strings.TrimSpace(
									v6.Data.Text[:min(60, len(v6.Data.Text))],
								),
								"\n",
								" | ",
							))
						}
					}
				}
			}
		}
	}
	log.Info("workspaces", "tot", wTot)
	log.Info("spaces", "tot", sTot)
	log.Info("folders", "tot", fTot)
	log.Info("lists", "tot", lTot)
	log.Info("tasks", "tot", tTot)
	log.Info("comments", "tot", cTot)
	return nil
}
