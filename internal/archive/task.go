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

func SaveTask(parent *List, t *api.Task, update bool) (*Task, error) {
	var tOld *Task = nil
	for _, c := range parent.Children {
		if c.Data.ID != t.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("task '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		tOld = c
		break
	}
	dir := taskDir(parent.GetDir(), t.ID)
	if err := jsonSave(t, dir); err != nil {
		return nil, err
	}
	if tOld == nil {
		log.Info("task created", "id", t.ID, "name", t.Name)
		tNew := &Task{Parent: parent, Data: t, Children: make([]*Comment, 0)}
		parent.Children = append(parent.Children, tNew)
		return tNew, nil
	}
	log.Warn("task updated", "id_old", tOld.Data.ID, "name_old", tOld.Data.Name, "id_new", t.ID, "name_new", t.Name)
	tOld.Data = t
	return tOld, nil
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
