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

func SaveList(parent *Folder, l *api.List, update bool) (*List, error) {
	var lOld *List = nil
	for _, c := range parent.Children {
		if c.Data.ID != l.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("list '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		lOld = c
		break
	}
	dir := listDir(parent.GetDir(), l.ID)
	if err := jsonSave(l, dir); err != nil {
		return nil, err
	}
	if lOld == nil {
		log.Info("list created", "id", l.ID, "name", l.Name)
		lNew := &List{Parent: parent, Data: l, Children: make([]*Task, 0)}
		parent.Children = append(parent.Children, lNew)
		return lNew, nil
	}
	log.Warn("list updated", "id_old", lOld.Data.ID, "name_old", lOld.Data.Name, "id_new", l.ID, "name_new", l.Name)
	lOld.Data = l
	return lOld, nil
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
