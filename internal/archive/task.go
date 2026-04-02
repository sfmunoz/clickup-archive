package archive

import (
	"encoding/json"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Task struct {
	Parent   *List
	Data     *api.Task
	Children []*Comment
}

func LoadTask(parent *List, id string) (*Task, error) {
	dir := taskDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.Task
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	t := &Task{
		Parent:   parent,
		Data:     &data,
		Children: make([]*Comment, 0),
	}
	entries, err := os.ReadDir(commentsDir(dir))
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
		c, err := LoadComment(t, e.Name())
		if err != nil {
			return nil, err
		}
		t.Children = append(t.Children, c)
	}
	return t, nil
}

func (t *Task) GetDir() string {
	return taskDir(t.Parent.GetDir(), t.Data.ID)
}
