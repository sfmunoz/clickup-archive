package archive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type List struct {
	Parent   *Folder
	Data     *api.List
	Children []*Task
}

func LoadList(parent *Folder, id string) (*List, error) {
	dir := listDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.List
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	l := &List{
		Parent:   parent,
		Data:     &data,
		Children: make([]*Task, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t, err := LoadTask(l, e.Name())
		if err != nil {
			return nil, err
		}
		l.Children = append(l.Children, t)
	}
	return l, nil
}

func (l *List) GetDir() string {
	return listDir(l.Parent.GetDir(), l.Data.ID)
}

func (l *List) SaveTask(t *api.Task, update bool) (*Task, error) {
	var tOld *Task = nil
	for _, c := range l.Children {
		if c.Data.ID != t.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("task '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		tOld = c
		break
	}
	dir := taskDir(l.GetDir(), t.ID)
	if err := jsonSave(t, dir); err != nil {
		return nil, err
	}
	if tOld == nil {
		log.Info("task created", "id", t.ID, "name", t.Name)
		tNew := &Task{Parent: l, Data: t, Children: make([]*Comment, 0)}
		l.Children = append(l.Children, tNew)
		return tNew, nil
	}
	log.Warn("task updated", "id_old", tOld.Data.ID, "name_old", tOld.Data.Name, "id_new", t.ID, "name_new", t.Name)
	tOld.Data = t
	return tOld, nil
}
