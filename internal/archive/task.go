package archive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Task struct {
	Parent   *List
	Data     api.Task
	Children []*Comment
}

func NewTask(parent *List, dir string) (*Task, error) {
	dir = filepath.Join(parent.GetDir(), dir)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filepath.Join(dir, "index.json"))
	if err != nil {
		return nil, err
	}
	var data api.Task
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	t := &Task{
		Parent:   parent,
		Data:     data,
		Children: make([]*Comment, 0),
	}
	commentsDir := filepath.Join(dir, "comments")
	entries, err := os.ReadDir(commentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return t, nil
		}
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		c, err := NewComment(t, e.Name())
		if err != nil {
			return nil, err
		}
		t.Children = append(t.Children, c)
	}
	return t, nil
}

func (t *Task) GetDir() string {
	return filepath.Join(t.Parent.GetDir(), t.Data.ID)
}
