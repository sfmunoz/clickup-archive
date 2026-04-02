package archive

import (
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Space struct {
	Parent   *Workspace
	Data     api.Space
	Children []*Folder
}

func NewSpace(parent *Workspace, dir string) (*Space, error) {
	dir = filepath.Join(parent.GetDir(), dir)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	s := &Space{
		Parent:   parent,
		Data:     api.Space{},
		Children: make([]*Folder, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		f, err := NewFolder(s, e.Name())
		if err != nil {
			return nil, err
		}
		s.Children = append(s.Children, f)
	}
	return s, nil
}

func (s *Space) GetDir() string {
	return filepath.Join(s.Parent.GetDir(), s.Data.ID)
}
