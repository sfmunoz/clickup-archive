package archive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type List struct {
	Parent   *Folder
	Data     *api.List
	Children []*Task
}

func LoadList(parent *Folder, dir string) (*List, error) {
	dir = filepath.Join(parent.GetDir(), dir)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filepath.Join(dir, "index.json"))
	if err != nil {
		return nil, err
	}
	var data api.List
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	l := &List{
		Parent:   parent,
		Data:     &data,
		Children: make([]*Task, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t, err := LoadTask(l, e.Name())
		if err != nil {
			return nil, err
		}
		l.Children = append(l.Children, t)
	}
	return l, nil
}

func (l *List) GetDir() string {
	return filepath.Join(l.Parent.GetDir(), l.Data.ID)
}
