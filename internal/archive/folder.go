package archive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Folder struct {
	Parent   *Space
	Data     *api.Folder
	Children []*List
}

func LoadFolder(parent *Space, id string) (*Folder, error) {
	dir := folderDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.Folder
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	f := &Folder{
		Parent:   parent,
		Data:     &data,
		Children: make([]*List, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		l, err := LoadList(f, e.Name())
		if err != nil {
			return nil, err
		}
		f.Children = append(f.Children, l)
	}
	return f, nil
}

func (f *Folder) GetDir() string {
	return folderDir(f.Parent.GetDir(), f.Data.ID)
}

func (f *Folder) SaveList(l *api.List, update bool) (*List, error) {
	var lOld *List = nil
	for _, c := range f.Children {
		if c.Data.ID != l.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("list '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		lOld = c
		break
	}
	dir := listDir(f.GetDir(), l.ID)
	if err := jsonSave(l, dir); err != nil {
		return nil, err
	}
	if lOld == nil {
		log.Info("list created", "id", l.ID, "name", l.Name)
		lNew := &List{Parent: f, Data: l, Children: make([]*Task, 0)}
		f.Children = append(f.Children, lNew)
		return lNew, nil
	}
	log.Warn("list updated", "id_old", lOld.Data.ID, "name_old", lOld.Data.Name, "id_new", l.ID, "name_new", l.Name)
	lOld.Data = l
	return lOld, nil
}
