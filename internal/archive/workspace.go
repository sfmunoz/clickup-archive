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

func (w *Workspace) GetDir() string {
	return workspaceDir(w.Parent.GetDir(), w.Data.ID)
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

func SaveWorkspace(parent *Archive, w *api.Workspace, update bool) (*Workspace, error) {
	var wOld *Workspace = nil
	for _, c := range parent.Children {
		if c.Data.ID != w.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("workspace '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		wOld = c
		break
	}
	dir := workspaceDir(parent.GetDir(), w.ID)
	if err := jsonSave(w, dir); err != nil {
		return nil, err
	}
	if wOld == nil {
		log.Info("workspace created", "id", w.ID, "name", w.Name)
		wNew := &Workspace{Parent: parent, Data: w, Children: make([]*Space, 0)}
		parent.Children = append(parent.Children, wNew)
		return wNew, nil
	}
	log.Warn("workspace updated", "id_old", wOld.Data.ID, "name_old", wOld.Data.Name, "id_new", w.ID, "name_new", w.Name)
	wOld.Data = w
	return wOld, nil
}
