package archive

import (
	"encoding/json"
	"fmt"
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

func (t *Task) SaveComment(c *api.Comment, update bool) error {
	var cOld *Comment = nil
	for _, ch := range t.Children {
		if ch.Data.ID != c.ID {
			continue
		}
		if !update {
			return fmt.Errorf("comment '%s' already exists and 'update' is false", ch.Data.ID)
		}
		cOld = ch
		break
	}
	if cOld == nil {
		log.Info("creating comment", "id", c.ID)
	} else {
		log.Warn("updating comment", "id_old", cOld.Data.ID, "id_new", c.ID)
	}
	dir := commentDir(t.GetDir(), c.ID)
	if err := jsonSave(c, dir); err != nil {
		return err
	}
	if cOld == nil {
		t.Children = append(t.Children, &Comment{Parent: t, Data: c, Children: make([]*struct{}, 0)})
	} else {
		cOld.Data = c
	}
	return nil
}

func (t *Task) IsCommentsDone() bool {
	_, err := os.Stat(doneFile(t.GetDir()))
	return err == nil
}

func (t *Task) MarkCommentsDone() error {
	return os.WriteFile(doneFile(t.GetDir()), []byte{}, 0o644)
}

func (t *Task) ClearComments() error {
	t.Children = make([]*Comment, 0)
	return os.RemoveAll(commentsDir(t.GetDir()))
}
