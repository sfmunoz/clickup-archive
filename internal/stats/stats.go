package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sfmunoz/clickup-archive/internal/archive"

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
	entries, err := os.ReadDir(s.clickupDir)
	if err != nil {
		return fmt.Errorf("read clickup dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if err := s.processWorkspace(filepath.Join(s.clickupDir, e.Name()), &counts); err != nil {
			return err
		}
	}
	log.Info("workspaces", "tot", counts.Workspaces)
	log.Info("spaces", "tot", counts.Spaces)
	log.Info("folders", "tot", counts.Folders)
	log.Info("lists", "tot", counts.Lists)
	log.Info("tasks", "tot", counts.Tasks)
	log.Info("comments", "tot", counts.Comments)
	a, err := archive.NewArchive(s.clickupDir)
	if err != nil {
		return err
	}
	for i1, v1 := range a.Children {
		log.Info("..", "i", i1, "v", v1.Data.ID, "n", v1.Data.Name)
		for i2, v2 := range v1.Children {
			log.Info("....", "i", i2, "v", v2.Data.ID, "n", v2.Data.Name)
			for i3, v3 := range v2.Children {
				log.Info("......", "i", i3, "v", v3.Data.ID, "n", v3.Data.Name)
				for i4, v4 := range v3.Children {
					log.Info("........", "i", i4, "v", v4.Data.ID, "n", v4.Data.Name)
					for i5, v5 := range v4.Children {
						log.Info("..........", "i", i5, "v", v5.Data.ID, "n", v5.Data.Name)
						for i6, v6 := range v5.Children {
							log.Info("............", "i", i6, "v", v6.Data.ID, "n", strings.ReplaceAll(
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
	return nil
}

func readIndex(dir string) (id, name string, err error) {
	data, err := os.ReadFile(filepath.Join(dir, "index.json"))
	if err != nil {
		return "", "", err
	}
	var entry struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", "", fmt.Errorf("parse index.json in %s: %w", dir, err)
	}
	return entry.ID, entry.Name, nil
}

func (s *Stats) processWorkspace(dir string, counts *Counts) error {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return nil
	}
	id, name, err := readIndex(dir)
	if err != nil {
		return err
	}
	counts.Workspaces++

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read workspace dir: %w", err)
	}

	var spaces, lists, tasks, comments int
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		sl, st, sc, err := s.processSpace(filepath.Join(dir, e.Name()), counts)
		if err != nil {
			return err
		}
		spaces++
		lists += sl
		tasks += st
		comments += sc
	}

	log.Info("workspace", "id", id, "name", name, "spaces", spaces, "lists", lists, "tasks", tasks, "comments", comments)
	return nil
}

func (s *Stats) processSpace(dir string, counts *Counts) (lists, tasks, comments int, err error) {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return 0, 0, 0, nil
	}
	id, name, err := readIndex(dir)
	if err != nil {
		return 0, 0, 0, err
	}
	counts.Spaces++

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read space dir: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		fl, ft, fc, ferr := s.processFolder(filepath.Join(dir, e.Name()), counts)
		if ferr != nil {
			return 0, 0, 0, ferr
		}
		lists += fl
		tasks += ft
		comments += fc
	}

	log.Info("  space", "id", id, "name", name, "lists", lists, "tasks", tasks, "comments", comments)
	return lists, tasks, comments, nil
}

func (s *Stats) processFolder(dir string, counts *Counts) (lists, tasks, comments int, err error) {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return 0, 0, 0, nil
	}
	counts.Folders++

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("read folder dir: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		lt, lc, lerr := s.processList(filepath.Join(dir, e.Name()), counts)
		if lerr != nil {
			return 0, 0, 0, lerr
		}
		lists++
		tasks += lt
		comments += lc
	}

	return lists, tasks, comments, nil
}

func (s *Stats) processList(dir string, counts *Counts) (tasks, comments int, err error) {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return 0, 0, nil
	}
	id, name, err := readIndex(dir)
	if err != nil {
		return 0, 0, err
	}
	counts.Lists++

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, 0, fmt.Errorf("read list dir: %w", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		tc, terr := s.processTask(filepath.Join(dir, e.Name()), counts)
		if terr != nil {
			return 0, 0, terr
		}
		tasks++
		comments += tc
	}

	log.Info("    list", "id", id, "name", name, "tasks", tasks, "comments", comments)
	return tasks, comments, nil
}

func (s *Stats) processTask(dir string, counts *Counts) (comments int, err error) {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return 0, nil
	}
	id, name, err := readIndex(dir)
	if err != nil {
		return 0, err
	}
	counts.Tasks++

	commentsDir := filepath.Join(dir, "comments")
	entries, err := os.ReadDir(commentsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, fmt.Errorf("read comments dir: %w", err)
		}
		log.Info("      task", "id", id, "name", name, "comments", 0)
		return 0, nil
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(commentsDir, e.Name(), "index.json")); os.IsNotExist(err) {
			continue
		}
		comments++
	}
	counts.Comments += comments

	log.Info("      task", "id", id, "name", name, "comments", comments)
	return comments, nil
}
