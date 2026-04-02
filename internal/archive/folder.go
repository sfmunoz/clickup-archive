package archive

import (
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Folder struct {
	Parent   *Space
	Data     api.Folder
	Children []*List
}

func NewFolder(parent *Space, dir string) (*Folder, error) {
	dir = filepath.Join(parent.GetDir(), dir)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	f := &Folder{
		Parent:   parent,
		Data:     api.Folder{},
		Children: make([]*List, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		l, err := NewList(f, e.Name())
		if err != nil {
			return nil, err
		}
		f.Children = append(f.Children, l)
	}
	return f, nil
}

func (f *Folder) GetDir() string {
	return filepath.Join(f.Parent.GetDir(), f.Data.ID)
}
