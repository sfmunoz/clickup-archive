package archive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Workspace struct {
	Parent   *Archive
	Data     *api.Workspace
	Children []*Space
}

func LoadWorkspace(parent *Archive, id string) (*Workspace, error) {
	dir := workspaceDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.Workspace
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	w := &Workspace{
		Parent:   parent,
		Data:     &data,
		Children: make([]*Space, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		s, err := LoadSpace(w, e.Name())
		if err != nil {
			return nil, err
		}
		w.Children = append(w.Children, s)
	}
	return w, nil
}

func (w *Workspace) GetDir() string {
	return workspaceDir(w.Parent.GetDir(), w.Data.ID)
}

func (w *Workspace) SaveSpace(s *api.Space, update bool) (*Space, error) {
	var sOld *Space = nil
	for _, c := range w.Children {
		if c.Data.ID != s.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("space '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		sOld = c
		break
	}
	dir := spaceDir(w.GetDir(), s.ID)
	if err := jsonSave(s, dir); err != nil {
		return nil, err
	}
	if sOld == nil {
		log.Info("space created", "id", s.ID, "name", s.Name)
		sNew := &Space{Parent: w, Data: s, Children: make([]*Folder, 0)}
		w.Children = append(w.Children, sNew)
		return sNew, nil
	}
	log.Warn("space updated", "id_old", sOld.Data.ID, "name_old", sOld.Data.Name, "id_new", s.ID, "name_new", s.Name)
	sOld.Data = s
	return sOld, nil
}
