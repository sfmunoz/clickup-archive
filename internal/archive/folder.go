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

func (f *Folder) GetDir() string {
	return folderDir(f.Parent.GetDir(), f.Data.ID)
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

func SaveFolder(parent *Space, f *api.Folder, update bool) (*Folder, error) {
	var fOld *Folder = nil
	for _, c := range parent.Children {
		if c.Data.ID != f.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("folder '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		fOld = c
		break
	}
	dir := folderDir(parent.GetDir(), f.ID)
	if err := jsonSave(f, dir); err != nil {
		return nil, err
	}
	if fOld == nil {
		log.Info("folder created", "id", f.ID, "name", f.Name)
		fNew := &Folder{Parent: parent, Data: f, Children: make([]*List, 0)}
		parent.Children = append(parent.Children, fNew)
		return fNew, nil
	}
	log.Warn("folder updated", "id_old", fOld.Data.ID, "name_old", fOld.Data.Name, "id_new", f.ID, "name_new", f.Name)
	fOld.Data = f
	return fOld, nil
}
