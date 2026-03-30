package stats

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Stats struct {
	clickupDir string
}

type Counts struct {
	Workspaces int
	Spaces     int
	Folders    int
	Lists      int
	Tasks      int
	Comments   int
}

func NewStats(clickupDir string) (*Stats, error) {
	return &Stats{clickupDir: clickupDir}, nil
}

func (s *Stats) Run() error {
	var counts Counts
	workspaceEntries, err := os.ReadDir(s.clickupDir)
	if err != nil {
		return fmt.Errorf("read clickup dir: %w", err)
	}
	for _, we := range workspaceEntries {
		if !we.IsDir() {
			continue
		}
		counts.Workspaces++
		wsDir := filepath.Join(s.clickupDir, we.Name())
		spaceEntries, err := os.ReadDir(wsDir)
		if err != nil {
			return fmt.Errorf("read workspace dir: %w", err)
		}
		for _, se := range spaceEntries {
			if !se.IsDir() {
				continue
			}
			counts.Spaces++
			spDir := filepath.Join(wsDir, se.Name())
			folderEntries, err := os.ReadDir(spDir)
			if err != nil {
				return fmt.Errorf("read space dir: %w", err)
			}
			for _, fe := range folderEntries {
				if !fe.IsDir() {
					continue
				}
				counts.Folders++
				foDir := filepath.Join(spDir, fe.Name())
				listEntries, err := os.ReadDir(foDir)
				if err != nil {
					return fmt.Errorf("read folder dir: %w", err)
				}
				for _, le := range listEntries {
					if !le.IsDir() {
						continue
					}
					counts.Lists++
					liDir := filepath.Join(foDir, le.Name())
					taskEntries, err := os.ReadDir(liDir)
					if err != nil {
						return fmt.Errorf("read list dir: %w", err)
					}
					for _, te := range taskEntries {
						if !te.IsDir() {
							continue
						}
						counts.Tasks++
						commentsDir := filepath.Join(liDir, te.Name(), "comments")
						commentEntries, err := os.ReadDir(commentsDir)
						if err != nil {
							if !os.IsNotExist(err) {
								return fmt.Errorf("read comments dir: %w", err)
							}
						} else {
							for _, ce := range commentEntries {
								if ce.IsDir() {
									counts.Comments++
								}
							}
						}
					}
				}
			}
		}
	}
	log.Info("workspaces", "tot", counts.Workspaces)
	log.Info("spaces", "tot", counts.Spaces)
	log.Info("folders", "tot", counts.Folders)
	log.Info("lists", "tot", counts.Lists)
	log.Info("tasks", "tot", counts.Tasks)
	log.Info("comments", "tot", counts.Comments)
	return nil
}
