package archive

import (
	"encoding/json"
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
