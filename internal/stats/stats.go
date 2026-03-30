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
	return nil
}

func (s *Stats) processWorkspace(dir string, counts *Counts) error {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return nil
	}
	counts.Workspaces++
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read workspace dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if err := s.processSpace(filepath.Join(dir, e.Name()), counts); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stats) processSpace(dir string, counts *Counts) error {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return nil
	}
	counts.Spaces++
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read space dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if err := s.processFolder(filepath.Join(dir, e.Name()), counts); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stats) processFolder(dir string, counts *Counts) error {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return nil
	}
	counts.Folders++
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read folder dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if err := s.processList(filepath.Join(dir, e.Name()), counts); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stats) processList(dir string, counts *Counts) error {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return nil
	}
	counts.Lists++
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read list dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if err := s.processTask(filepath.Join(dir, e.Name()), counts); err != nil {
			return err
		}
	}
	return nil
}

func (s *Stats) processTask(dir string, counts *Counts) error {
	if _, err := os.Stat(filepath.Join(dir, "index.json")); os.IsNotExist(err) {
		return nil
	}
	counts.Tasks++
	commentsDir := filepath.Join(dir, "comments")
	entries, err := os.ReadDir(commentsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("read comments dir: %w", err)
		}
		return nil
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(commentsDir, e.Name(), "index.json")); os.IsNotExist(err) {
			continue
		}
		counts.Comments++
	}
	return nil
}
