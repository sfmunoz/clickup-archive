package archive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Workspace struct {
	Parent   *Archive
	Data     *api.Workspace
	Children []*Space
}

func NewWorkspace(parent *Archive, dir string) (*Workspace, error) {
	dir = filepath.Join(parent.GetDir(), dir)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filepath.Join(dir, "index.json"))
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
		s, err := NewSpace(w, e.Name())
		if err != nil {
			return nil, err
		}
		w.Children = append(w.Children, s)
	}
	return w, nil
}

func (w *Workspace) GetDir() string {
	return filepath.Join(w.Parent.GetDir(), w.Data.ID)
}
