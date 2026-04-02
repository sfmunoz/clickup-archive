package archive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Folder struct {
	Parent   *Space
	Data     *api.Folder
	Children []*List
}

func LoadFolder(parent *Space, id string) (*Folder, error) {
	dir := filepath.Join(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filepath.Join(dir, "index.json"))
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
	return filepath.Join(f.Parent.GetDir(), f.Data.ID)
}
