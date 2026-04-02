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

func SaveSpace(parent *Workspace, s *api.Space, update bool) (*Space, error) {
	var sOld *Space = nil
	for _, c := range parent.Children {
		if c.Data.ID != s.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("space '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		sOld = c
		break
	}
	dir := spaceDir(parent.GetDir(), s.ID)
	if err := jsonSave(s, dir); err != nil {
		return nil, err
	}
	if sOld == nil {
		log.Info("space created", "id", s.ID, "name", s.Name)
		sNew := &Space{Parent: parent, Data: s, Children: make([]*Folder, 0)}
		parent.Children = append(parent.Children, sNew)
		return sNew, nil
	}
	log.Warn("space updated", "id_old", sOld.Data.ID, "name_old", sOld.Data.Name, "id_new", s.ID, "name_new", s.Name)
	sOld.Data = s
	return sOld, nil
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
