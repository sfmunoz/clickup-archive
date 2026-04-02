package archive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Space struct {
	Parent   *Workspace
	Data     *api.Space
	Children []*Folder
}

func LoadSpace(parent *Workspace, id string) (*Space, error) {
	dir := spaceDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.Space
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	s := &Space{
		Parent:   parent,
		Data:     &data,
		Children: make([]*Folder, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		f, err := LoadFolder(s, e.Name())
		if err != nil {
			return nil, err
		}
		s.Children = append(s.Children, f)
	}
	return s, nil
}

func (s *Space) GetDir() string {
	return spaceDir(s.Parent.GetDir(), s.Data.ID)
}

func (s *Space) SaveFolder(f *api.Folder, update bool) (*Folder, error) {
	var fOld *Folder = nil
	for _, c := range s.Children {
		if c.Data.ID != f.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("folder '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		fOld = c
		break
	}
	dir := folderDir(s.GetDir(), f.ID)
	if err := jsonSave(f, dir); err != nil {
		return nil, err
	}
	if fOld == nil {
		log.Info("folder created", "id", f.ID, "name", f.Name)
		fNew := &Folder{Parent: s, Data: f, Children: make([]*List, 0)}
		s.Children = append(s.Children, fNew)
		return fNew, nil
	}
	log.Warn("folder updated", "id_old", fOld.Data.ID, "name_old", fOld.Data.Name, "id_new", f.ID, "name_new", f.Name)
	fOld.Data = f
	return fOld, nil
}
